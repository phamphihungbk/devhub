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
	PluginID         string                 `json:"plugin_id" validate:"required,uuid,enabled_plugin"`
	ProjectID        string                 `json:"project_id" validate:"required,uuid"`
	EnvironmentID    string                 `json:"environment_id" validate:"required,uuid"`
	RequestedBy      string                 `json:"requested_by" validate:"required,uuid"`
	ApprovalResource string                 `json:"approval_resource" validate:"omitempty"`
	Variables        map[string]interface{} `json:"variables" validate:"required,plugin_config_schema"`
}

func (u *scaffoldRequestUsecase) CreateScaffoldRequest(ctx context.Context, input CreateScaffoldRequestInput) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[usecase scaffold_request/create_scaffold_request CreateScaffoldRequest] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
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

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	projectID := uuid.MustParse(input.ProjectID)
	pluginID := uuid.MustParse(input.PluginID)
	requestedBy := uuid.MustParse(input.RequestedBy)

	approvalResource, err := new(entity.ApprovalResource).Parse(input.ApprovalResource)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid approval resource", nil))
	}

	scaffoldRequest = &entity.ScaffoldRequest{
		ProjectID:   projectID,
		RequestedBy: requestedBy,
		PluginID:    pluginID,
		Status:      entity.ScaffoldRequestPending,
		Variables:   entity.ScaffoldRequestVariables(input.Variables),
	}
	created, err := u.scaffoldRequestRepository.CreateOne(ctx, scaffoldRequest)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create scaffold request", nil))
	}

	policy, err := u.approvalRepository.FindMatchingApprovalPolicy(ctx, repository.FindMatchingApprovalPolicyInput{
		Resource:      approvalResource.String(),
		EnvironmentID: uuid.MustParse(input.EnvironmentID),
	})
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find approval policy", nil))
	}

	if policy != nil {
		if _, err := u.approvalRepository.CreateApprovalRequest(ctx, &entity.ApprovalRequest{
			Resource:          policy.Resource,
			ResourceID:        created.ID,
			RequestedBy:       requestedBy,
			Status:            entity.ApprovalRequestStatusPending,
			RequiredApprovals: policy.RequiredApprovals,
			ApprovedCount:     0,
			RejectedCount:     0,
		}); err != nil {
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create approval request", nil))
		}
	}

	return created, nil
}
