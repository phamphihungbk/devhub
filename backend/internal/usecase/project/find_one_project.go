package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"
	"errors"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type FindOneProjectInput struct {
	ID string `json:"id" validate:"required,uuid4"`
}

type ProjectDetail struct {
	ID            uuid.UUID
	Name          string
	Description   string
	OwnerTeamName string
	CreatorName   string
}

func (u *projectUsecase) FindOneProject(ctx context.Context, input FindOneProjectInput) (projectDetail *ProjectDetail, err error) {
	const errLocation = "[usecase project/find_one_project FindOneProject] "
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

	project, err := u.projectRepository.FindOne(ctx, userID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find project by ID", nil))
		}
		return nil, err // Return the NotFoundError directly
	}

	return u.enrichProjectDetail(ctx, project), nil
}

func (u *projectUsecase) enrichProjectDetail(ctx context.Context, project *entity.Project) *ProjectDetail {
	if project == nil {
		return nil
	}

	detail := &ProjectDetail{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
	}

	// TODO: maybe should move to repo with join query instead
	ownerTeam, err := u.teamRepository.FindOne(ctx, project.OwnerTeamID)
	if err == nil && ownerTeam != nil {
		detail.OwnerTeamName = ownerTeam.Name
	}

	creator, err := u.userRepository.FindOne(ctx, project.CreatedBy)
	if err == nil && creator != nil {
		detail.CreatorName = creator.Name
	}

	return detail
}
