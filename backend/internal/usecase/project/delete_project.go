package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type DeleteProjectInput struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (u *projectUsecase) DeleteProject(ctx context.Context, input DeleteProjectInput) (project *entity.Project, err error) {
	const errLocation = "[usecase project/delete_project DeleteProject] "
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
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid project ID", nil))
	}

	deleted, err := u.projectRepository.DeleteOne(ctx, userID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete project", nil))
	}

	return deleted, nil
}
