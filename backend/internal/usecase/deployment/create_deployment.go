package usecase

import (
	"context"

	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type CreateDeploymentInput struct {
	ProjectID   string `json:"project_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
	Environment string `json:"environment" example:"prod" validate:"required,oneof=dev staging prod"`
	Service     string `json:"service" example:"Service Name" validate:"required"`
	Version     string `json:"version" example:"v1.0.0" validate:"required,git_revision"`
	TriggeredBy string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required,uuid"`
}

func (u *deploymentUsecase) CreateDeployment(ctx context.Context, input CreateDeploymentInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[usecase deployment/create_deployment CreateDeployment] "
	defer misc.WrapErrorWithPrefix(errLocation, &err)

	// Create a new validator instance
	vInstance, err := validator.NewValidator(
		validator.WithTagNameFunc(validator.JSONTagNameFunc),
		validator.WithCustomValidator(validator.GitRevisionValidator{}),
	)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create validator", nil))
	}

	// Validate Input
	err = vInstance.Struct(input)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
	}

	env, err := new(entity.ProjectEnvironment).Parse(input.Environment)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid environment", nil))
	}

	deployment = &entity.Deployment{
		ProjectID:   uuid.MustParse(input.ProjectID),
		Environment: env,
		Service:     input.Service,
		Version:     input.Version,
		Status:      entity.StatusPending,
		TriggeredBy: uuid.MustParse(input.TriggeredBy),
	}

	created, err := u.deploymentRepository.CreateOne(ctx, deployment)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create deployment", nil))
	}

	return created, nil
}
