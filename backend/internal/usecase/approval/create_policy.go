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
	Resource          string  `json:"resource" validate:"required,oneof=scaffold_request deployment"`
	ServiceID         *string `json:"service_id" validate:"omitempty,uuid"`
	EnvironmentID     string  `json:"environment_id" validate:"required,uuid"`
	RequiredApprovals int     `json:"required_approvals" validate:"required,min=1"`
	Enabled           bool    `json:"enabled"`
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

	var serviceID *uuid.UUID
	if input.ServiceID != nil {
		value := uuid.MustParse(*input.ServiceID)
		serviceID = &value
	}
	environmentID := uuid.MustParse(input.EnvironmentID)

	policy := &entity.ApprovalPolicy{
		Resource:          input.Resource,
		ServiceID:         serviceID,
		EnvironmentID:     environmentID,
		RequiredApprovals: input.RequiredApprovals,
		Enabled:           input.Enabled,
	}

	created, err := u.approvalRepository.CreateApprovalPolicy(ctx, policy)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create approval policy", nil))
	}

	return created, nil
}
