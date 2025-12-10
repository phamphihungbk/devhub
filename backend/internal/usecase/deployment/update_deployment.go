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
	ID          string  `json:"id" validate:"required,uuid"`
	Environment *string `json:"environment" validate:"required,min=2,max=100"`
	Service     *string `json:"service" validate:"required,min=2,max=100"`
	Version     *string `json:"version" validate:"required,min=1,max=50"`
	Status      *string `json:"status" validate:"required,min=2,max=100"`
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
		ID:          uuid.MustParse(input.ID),
		Environment: (*entity.ProjectEnvironment)(input.Environment),
		Status:      (*entity.DeploymentStatus)(input.Status),
		Service:     input.Service,
		Version:     input.Version,
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update deployment", nil))
	}

	return updated, nil
}
