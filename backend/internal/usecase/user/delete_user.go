package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type DeleteUserInput struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (u *userUsecase) DeleteUser(ctx context.Context, input DeleteUserInput) (user *entity.User, err error) {
	const errLocation = "[usecase user/delete_user DeleteUser] "
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

	userID, err := uuid.Parse(input.ID)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid user ID", nil))
	}

	deleted, err := u.userRepository.DeleteOne(ctx, userID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete user", nil))
	}

	return deleted, nil
}
