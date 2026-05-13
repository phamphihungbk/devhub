package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateDeploymentInput struct {
	ServiceID        string `json:"service_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
	PluginID         string `json:"plugin_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
	ReleaseID        string `json:"release_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"omitempty,uuid"`
	EnvironmentID    string `json:"environment_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
	TriggeredBy      string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
	ApprovalResource string `json:"approval_resource" validate:"omitempty"`
	ApprovalAction   string `json:"approval_action" validate:"omitempty"`
}

func (u *deploymentUsecase) CreateDeployment(ctx context.Context, input CreateDeploymentInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[usecase deployment/create_deployment CreateDeployment] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
		validator.WithCustomValidator(validator.EnabledPluginValidator{
			Context:          ctx,
			PluginRepository: u.pluginRepository,
		}),
		validator.WithCustomValidator(&validator.PluginConfigSchemaValidator{
			Context:          ctx,
			PluginRepository: u.pluginRepository,
		}),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	serviceID := uuid.MustParse(input.ServiceID)
	environmentID := uuid.MustParse(input.EnvironmentID)
	pluginID := uuid.MustParse(input.PluginID)
	triggeredBy := uuid.MustParse(input.TriggeredBy)

	plugin, err := u.pluginRepository.FindOne(ctx, pluginID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid deployment plugin", nil))
	}
	if !plugin.Enabled {
		return nil, errs.NewBadRequestError("deployment plugin is disabled", nil)
	}

	release, err := u.resolveDeploymentRelease(ctx, serviceID, input.ReleaseID)
	if err != nil {
		return nil, err
	}

	deployment = &entity.Deployment{
		ServiceID:     serviceID,
		EnvironmentID: environmentID,
		ReleaseID:     release.ID,
		PluginID:      pluginID,
		Status:        entity.DeploymentStatusPending,
		TriggeredBy:   triggeredBy,
	}

	created, err := u.deploymentRepository.CreateOne(ctx, deployment)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create deployment", nil))
	}

	hasApproval, err := u.createDeploymentApprovalRequest(ctx, input, created.ID, serviceID, environmentID, triggeredBy)
	if err != nil {
		return nil, err
	}
	if !hasApproval {
		if err := u.createDeploymentJob(ctx, created, release, triggeredBy); err != nil {
			return nil, err
		}
	}

	return created, nil
}

func (u *deploymentUsecase) resolveDeploymentRelease(ctx context.Context, serviceID uuid.UUID, releaseID string) (*entity.Release, error) {
	if releaseID != "" {
		release, err := u.releaseRepository.FindOne(ctx, uuid.MustParse(releaseID))
		if err != nil {
			return nil, misc.WrapError(err, errs.NewBadRequestError("invalid release", nil))
		}
		if release.ServiceID != serviceID {
			return nil, errs.NewBadRequestError("release does not belong to service", nil)
		}
		return release, nil
	}

	releases, err := u.releaseRepository.FindAll(ctx, repository.FindAllReleasesFilter{ServiceID: serviceID})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to fetch releases", nil))
	}
	for _, release := range misc.GetValue(releases) {
		if release.Status == entity.ReleaseStatusCompleted {
			return &release, nil
		}
	}
	return nil, errs.NewBadRequestError("release_id is required when service has no completed release", nil)
}

func (u *deploymentUsecase) createDeploymentApprovalRequest(
	ctx context.Context,
	input CreateDeploymentInput,
	deploymentID uuid.UUID,
	serviceID uuid.UUID,
	environmentID uuid.UUID,
	triggeredBy uuid.UUID,
) (bool, error) {
	if input.ApprovalResource == "" || input.ApprovalAction == "" {
		return false, nil
	}

	approvalResource, err := new(entity.ApprovalResource).Parse(input.ApprovalResource)
	if err != nil {
		return false, misc.WrapError(err, errs.NewBadRequestError("invalid approval resource", nil))
	}
	approvalAction, err := new(entity.ApprovalAction).Parse(input.ApprovalAction)
	if err != nil {
		return false, misc.WrapError(err, errs.NewBadRequestError("invalid approval action", nil))
	}

	policy, err := u.approvalRepository.FindMatchingApprovalPolicy(ctx, repository.FindMatchingApprovalPolicyInput{
		Resource:      approvalResource.String(),
		EnvironmentID: environmentID,
		ServiceID:     &serviceID,
	})
	if err != nil {
		return false, misc.WrapError(err, errs.NewInternalServerError("failed to find approval policy", nil))
	}
	if policy == nil {
		return false, nil
	}

	if _, err := u.approvalRepository.CreateApprovalRequest(ctx, &entity.ApprovalRequest{
		Resource:          policy.Resource,
		Action:            approvalAction.String(),
		ResourceID:        deploymentID,
		RequestedBy:       triggeredBy,
		Status:            entity.ApprovalRequestStatusPending,
		RequiredApprovals: policy.RequiredApprovals,
		ApprovedCount:     0,
		RejectedCount:     0,
	}); err != nil {
		return false, misc.WrapError(err, errs.NewInternalServerError("failed to create approval request", nil))
	}

	return true, nil
}

func (u *deploymentUsecase) createDeploymentJob(ctx context.Context, deployment *entity.Deployment, release *entity.Release, triggeredBy uuid.UUID) error {
	variables, err := u.deploymentJobVariables(ctx, deployment, release)
	if err != nil {
		return err
	}
	if _, err := u.jobRepository.CreateOne(ctx, &entity.Job{
		Type:         entity.JobTypeDeployment,
		Status:       entity.JobStatusQueued,
		ResourceType: entity.JobResourceTypeDeployment,
		ResourceID:   deployment.ID,
		PluginID:     deployment.PluginID,
		Payload: entity.JobPayload{
			Action:    entity.JobTypeDeployment.String(),
			Variables: variables,
		},
		CreatedBy: triggeredBy,
	}); err != nil {
		return misc.WrapError(err, errs.NewInternalServerError("failed to create deployment job", nil))
	}
	return nil
}

func (u *deploymentUsecase) deploymentJobVariables(ctx context.Context, deployment *entity.Deployment, release *entity.Release) (map[string]interface{}, error) {
	service, err := u.serviceRepository.FindOne(ctx, deployment.ServiceID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find deployment service", nil))
	}
	environment, err := u.environmentRepository.FindOne(ctx, deployment.EnvironmentID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find deployment environment", nil))
	}
	return map[string]interface{}{
		"deployment_id": deployment.ID.String(),
		"project_id":    service.ProjectID.String(),
		"service_id":    service.ID.String(),
		"plugin_id":     deployment.PluginID.String(),
		"service":       service.Name,
		"environment":   environment.Name,
		"version":       release.Tag,
		"repo_url":      service.RepoURL,
	}, nil
}
