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
	Status       *string   `json:"status" validate:"required,oneof=draft active archived deprecated"`
	OwnerTeam    *string   `json:"owner_team" validate:"required,min=1,max=255"`
	RepoURL      *string   `json:"repo_url" validate:"required,max=2048"`
	RepoProvider *string   `json:"repo_provider" validate:"required,min=1,max=32"`
	OwnerContact *string   `json:"owner_contact" validate:"required,min=1,max=255"`
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

	var status *entity.ProjectStatus
	if input.Status != nil {
		parsedStatus, parseErr := new(entity.ProjectStatus).Parse(*input.Status)
		if parseErr != nil {
			return nil, misc.WrapError(parseErr, errs.NewBadRequestError("invalid project status", map[string]string{"details": parseErr.Error()}))
		}
		status = &parsedStatus
	}

	updated, err := u.projectRepository.UpdateOne(ctx, repository.UpdateProjectInput{
		ID:           uuid.MustParse(input.ID),
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
		Status:       status,
		OwnerTeam:    input.OwnerTeam,
		RepoURL:      input.RepoURL,
		RepoProvider: input.RepoProvider,
		OwnerContact: input.OwnerContact,
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update project", nil))
	}

	return updated, nil
}
