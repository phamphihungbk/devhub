package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"
)

type CreateDeploymentInput struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin user"`
}

func (u *deploymentUsecase) CreateDeployment(ctx context.Context, input CreateDeploymentInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[usecase deployment/create_deployment CreateDeployment] "
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

	deployment = &entity.Deployment{
		Name:  input.Name,
		Email: input.Email,
		Role:  entity.UserRole(input.Role),
	}

	created, err := u.deploymentRepository.CreateOne(ctx, deployment)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create deployment", nil))
	}

	return created, nil
}
