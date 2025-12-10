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

type UpdateProjectInput struct {
	ID           string    `json:"id" validate:"required,uuid"`
	Name         *string   `json:"name" validate:"min=0,max=100"`
	Description  *string   `json:"description" validate:"min=0,max=100"`
	Environments *[]string `json:"environments" validate:"dive,required"`
}

func (u *projectUsecase) UpdateProject(ctx context.Context, input UpdateProjectInput) (project *entity.Project, err error) {
	const errLocation = "[usecase project/update_project UpdateProject] "
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

	updated, err := u.projectRepository.UpdateOne(ctx, repository.UpdateProjectInput{
		ID:           uuid.MustParse(input.ID),
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update project", nil))
	}

	return updated, nil
}
