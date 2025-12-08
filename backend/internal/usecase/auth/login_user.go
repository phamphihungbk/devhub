package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"
)

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *authUsecase) LoginUser(ctx context.Context, input LoginUserInput) (user *entity.User, err error) {
	const errLocation = "[usecase auth/login_user LoginUser] "
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

	// deployment = &entity.Deployment{
	// 	Name:  input.Name,
	// 	Email: input.Email,
	// 	Role:  entity.UserRole(input.Role),
	// }

	// created, err := u.deploymentRepository.CreateOne(ctx, deployment)
	// if err != nil {
	// 	return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create deployment", nil))
	// }

	return nil, nil
}
