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

type CreateScaffoldRequestInput struct {
	PluginID         string                   `json:"plugin_id" validate:"required,uuid"`
	ProjectID        string                   `json:"project_id" validate:"required,uuid"`
	RequestedBy      string                   `json:"requested_by" validate:"required,uuid"`
	Environment      string                   `json:"environment" validate:"required,oneof=dev staging prod"`
	ApprovalResource string                   `json:"approval_resource" validate:"omitempty"`
	ApprovalAction   string                   `json:"approval_action" validate:"omitempty"`
	Variables        ScaffoldRequestVariables `json:"variables" validate:"required"`
}

type ScaffoldRequestVariables struct {
	ServiceName   string `json:"service_name" validate:"required"`
	ModulePath    string `json:"module_path" validate:"required"`
	Port          int    `json:"port" validate:"required"`
	Database      string `json:"database" validate:"required"`
	EnableLogging bool   `json:"enable_logging" validate:"required"`
}

func (u *scaffoldRequestUsecase) CreateScaffoldRequest(ctx context.Context, input CreateScaffoldRequestInput) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[usecase scaffold_request/create_scaffold_request CreateScaffoldRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	projectID, err := uuid.Parse(input.ProjectID)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid project ID", nil))
	}

	pluginID, err := uuid.Parse(input.PluginID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid plugin ID", nil))
	}

	requestedBy, err := uuid.Parse(input.RequestedBy)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid requested by user ID", nil))
	}

	approvalResource, err := new(entity.ApprovalResource).Parse(input.ApprovalResource)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid approval resource", nil))
	}

	approvalAction, err := new(entity.ApprovalAction).Parse(input.ApprovalAction)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid approval action", nil))
	}

	scaffoldRequest = &entity.ScaffoldRequest{
		PluginID:    pluginID,
		ProjectID:   projectID,
		RequestedBy: requestedBy,
		Status:      entity.ScaffoldRequestPending,
		Environment: new(entity.ProjectEnvironment).MustParse(input.Environment),
		Variables:   entity.ScaffoldRequestVariables(input.Variables),
	}
	created, err := u.scaffoldRequestRepository.CreateOne(ctx, scaffoldRequest)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create scaffold request", nil))
	}

	policy, err := u.approvalRepository.FindMatchingApprovalPolicy(ctx, repository.FindMatchingApprovalPolicyInput{
		Resource:    approvalResource.String(),
		Action:      approvalAction.String(),
		ProjectID:   misc.ToPointer(projectID),
		Environment: misc.ToPointer(input.Environment),
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find approval policy", nil))
	}

	if _, err := u.approvalRepository.CreateApprovalRequest(ctx, &entity.ApprovalRequest{
		Resource:          policy.Resource,
		Action:            policy.Action,
		ResourceID:        created.ID,
		RequestedBy:       requestedBy,
		ProjectID:         misc.ToPointer(projectID),
		Environment:       misc.ToPointer(input.Environment),
		Status:            entity.ApprovalRequestStatusPending,
		RequiredApprovals: policy.RequiredApprovals,
		ApprovedCount:     0,
		RejectedCount:     0,
	}); err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create approval request", nil))
	}

	return created, nil
}
