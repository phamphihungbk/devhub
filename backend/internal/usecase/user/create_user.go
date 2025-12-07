package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"
)

type CreateUserInput struct {
	Name  string `json:"name" validate:"required,gt=0"`
	Email string `json:"email" validate:"required,gt=0"`
	Role  string `json:"role" validate:"required,gt=0"`
}

func (u *userUsecase) CreateUser(ctx context.Context, input CreateUserInput) (user *entity.User, err error) {
	const errLocation = "[usecase user/create_user CreateUser] "
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

	user = &entity.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  entity.UserRole(input.Role),
	}

	created, err := u.userRepository.CreateOne(ctx, user)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create user", nil))
	}

	return created, nil
}
