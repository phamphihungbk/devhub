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
	Environment      string `json:"environment" example:"prod" validate:"required,oneof=dev staging prod"`
	Version          string `json:"version" example:"v1.0.0" validate:"required,git_revision,startswith=v"`
	TriggeredBy      string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
	ApprovalResource string `json:"approval_resource" validate:"omitempty"`
	ApprovalAction   string `json:"approval_action" validate:"omitempty"`
}

func (u *deploymentUsecase) CreateDeployment(ctx context.Context, input CreateDeploymentInput) (deployment *entity.Deployment, err error) {
	// TODO: later deployment rely on release instead
	const errLocation = "[usecase deployment/create_deployment CreateDeployment] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
		validator.WithCustomValidator(validator.GitRevisionValidator{}),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	env, err := new(entity.ProjectEnvironment).Parse(input.Environment)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid environment", nil))
	}

	serviceID := uuid.MustParse(input.ServiceID)
	pluginID := uuid.MustParse(input.PluginID)
	triggeredBy := uuid.MustParse(input.TriggeredBy)
	approvalResource, err := entity.ParseOptionalApprovalResource(input.ApprovalResource)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid approval resource", nil))
	}

	approvalAction, err := entity.ParseOptionalApprovalAction(input.ApprovalAction)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid approval action", nil))
	}

	plugin, err := u.pluginRepository.FindOne(ctx, pluginID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid deployment plugin", nil))
	}

	if plugin.Type != entity.PluginDeployer {
		return nil, errs.NewBadRequestError("deployment plugin must be a deployer", nil)
	}

	if !plugin.Enabled {
		return nil, errs.NewBadRequestError("deployment plugin is disabled", nil)
	}

	deployment = &entity.Deployment{
		ServiceID:   serviceID,
		PluginID:    pluginID,
		Environment: env,
		Version:     input.Version,
		Status:      entity.DeploymentStatusPending,
		TriggeredBy: triggeredBy,
	}

	created, err := u.deploymentRepository.CreateOne(ctx, deployment)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create deployment", nil))
	}

	if approvalResource != nil && approvalAction != nil {
		environment := input.Environment
		policy, err := u.approvalRepository.FindMatchingApprovalPolicy(ctx, repository.FindMatchingApprovalPolicyInput{
			Resource:    approvalResource.String(),
			Action:      approvalAction.String(),
			ServiceID:   &serviceID,
			Environment: &environment,
		})
		if err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find approval policy", nil))
		}

		if policy != nil {
			if _, err := u.approvalRepository.CreateApprovalRequest(ctx, &entity.ApprovalRequest{
				Resource:          policy.Resource,
				Action:            policy.Action,
				ResourceID:        created.ID,
				RequestedBy:       triggeredBy,
				ServiceID:         &serviceID,
				Environment:       &environment,
				Status:            entity.ApprovalRequestStatusPending,
				RequiredApprovals: policy.RequiredApprovals,
				ApprovedCount:     0,
				RejectedCount:     0,
			}); err != nil {
				return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create approval request", nil))
			}
		}
	}

	return created, nil
}
