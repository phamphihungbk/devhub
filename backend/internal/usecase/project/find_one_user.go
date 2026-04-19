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

func (u *projectUsecase) FindOneProject(ctx context.Context, input FindOneProjectInput) (project *entity.Project, err error) {
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

	project, err = u.projectRepository.FindOne(ctx, userID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find project by ID", nil))
		}
		return nil, err // Return the NotFoundError directly
	}

	u.enrichProjectCreator(ctx, project)

	return project, nil
}
