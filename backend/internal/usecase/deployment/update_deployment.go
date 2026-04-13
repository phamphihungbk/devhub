package usecase

import (
	"context"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type UpdateDeploymentInput struct {
	ID          string  `json:"id" validate:"required,uuid"`
	Environment *string `json:"environment" validate:"omitempty,oneof=dev staging prod"`
	Service     *string `json:"service" validate:"omitempty,min=2,max=100"`
	Version     *string `json:"version" validate:"omitempty,min=1,max=50"`
	Status      *string `json:"status" validate:"omitempty,oneof=pending running completed failed rolled_back"`
	ExternalRef *string `json:"external_ref" validate:"omitempty,max=255"`
	CommitSHA   *string `json:"commit_sha" validate:"omitempty,max=64"`
	FinishedAt  *string `json:"finished_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
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

	var finishedAt *time.Time
	if input.FinishedAt != nil {
		parsed, parseErr := time.Parse(time.RFC3339, *input.FinishedAt)
		if parseErr != nil {
			return nil, misc.WrapError(parseErr, errs.NewBadRequestError("invalid finished_at", map[string]string{"details": parseErr.Error()}))
		}
		finishedAt = &parsed
	}

	updated, err := u.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:          uuid.MustParse(input.ID),
		Environment: (*entity.ProjectEnvironment)(input.Environment),
		Status:      (*entity.DeploymentStatus)(input.Status),
		Service:     input.Service,
		Version:     input.Version,
		ExternalRef: input.ExternalRef,
		CommitSHA:   input.CommitSHA,
		FinishedAt:  finishedAt,
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update deployment", nil))
	}

	return updated, nil
}
