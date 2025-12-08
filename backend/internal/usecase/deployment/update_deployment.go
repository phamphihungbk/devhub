package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type UpdateDeploymentInput struct {
	ID   string  `json:"id" validate:"required,uuid"`
	Name *string `json:"name" validate:"required,min=2,max=100"`
	Role *string `json:"role" validate:"required,oneof=admin user"`
}

func (u *deploymentUsecase) UpdateDeployment(ctx context.Context, input UpdateDeploymentInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[usecase deployment/update_deployment UpdateDeployment] "
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

	updated, err := u.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:   uuid.MustParse(input.ID),
		Name: input.Name,
		Role: (*entity.UserRole)(input.Role),
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update deployment", nil))
	}

	return updated, nil
}
