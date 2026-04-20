package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateApprovalPolicyInput struct {
	Resource          string  `json:"resource" validate:"required"`
	Action            string  `json:"action" validate:"required"`
	ProjectID         *string `json:"project_id" validate:"omitempty,uuid"`
	ServiceID         *string `json:"service_id" validate:"omitempty,uuid"`
	Environment       *string `json:"environment" validate:"omitempty,oneof=dev staging prod"`
	RequiredApprovals int     `json:"required_approvals" validate:"required,min=1"`
	Enabled           *bool   `json:"enabled"`
}

func (u *approvalUsecase) CreateApprovalPolicy(ctx context.Context, input CreateApprovalPolicyInput) (_ *entity.ApprovalPolicy, err error) {
	const errLocation = "[usecase approval/create_policy CreateApprovalPolicy] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
	)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	if err := vInstance.Struct(input); err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	var (
		projectID *uuid.UUID
		serviceID *uuid.UUID
	)

	if input.ProjectID != nil {
		value := uuid.MustParse(*input.ProjectID)
		projectID = &value
	}
	if input.ServiceID != nil {
		value := uuid.MustParse(*input.ServiceID)
		serviceID = &value
	}

	policy := &entity.ApprovalPolicy{
		Resource:          input.Resource,
		Action:            input.Action,
		ProjectID:         projectID,
		ServiceID:         serviceID,
		Environment:       input.Environment,
		RequiredApprovals: input.RequiredApprovals,
		Enabled:           misc.GetValue(input.Enabled),
	}
	if input.Enabled == nil {
		policy.Enabled = true
	}

	created, err := u.approvalRepository.CreateApprovalPolicy(ctx, policy)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create approval policy", nil))
	}

	return created, nil
}
