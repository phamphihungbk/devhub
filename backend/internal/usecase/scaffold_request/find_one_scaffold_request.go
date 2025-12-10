package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"errors"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type FindOneScaffoldRequestInput struct {
	ID string `json:"id" validate:"required,uuid4"`
}

func (u *scaffoldRequestUsecase) FindOneScaffoldRequest(ctx context.Context, input FindOneScaffoldRequestInput) (scaffoldRequest *entity.ScaffoldRequest, err error) {
	const errLocation = "[usecase scaffold_request/find_one_scaffold_request FindOneScaffoldRequest] "
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

	scaffoldRequest, err = u.scaffoldRequestRepository.FindOne(ctx, scaffoldRequestID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find scaffold request by ID", nil))
		}
		return nil, err // Return the NotFoundError directly
	}

	return scaffoldRequest, nil
}
