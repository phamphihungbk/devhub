package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type ScaffoldRequestDeleteInput struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (u *scaffoldRequestUsecase) DeleteScaffoldRequest(ctx context.Context, input ScaffoldRequestDeleteInput) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[usecase scaffold_request/delete_scaffold_request DeleteScaffoldRequest] "
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

	scaffoldRequestID, err := uuid.Parse(input.ID)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid scaffold request ID", nil))
	}

	deleted, err := u.scaffoldRequestRepository.DeleteOne(ctx, scaffoldRequestID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete scaffold request", nil))
	}

	return deleted, nil
}
