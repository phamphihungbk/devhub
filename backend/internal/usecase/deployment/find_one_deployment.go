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

type FindOneDeploymentInput struct {
	ID string `json:"id" validate:"required,uuid4"`
}

func (u *deploymentUsecase) FindOneDeployment(ctx context.Context, input FindOneDeploymentInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[usecase deployment/find_one_deployment FindOneDeployment] "
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

	deploymentID, err := uuid.Parse(input.ID)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid deployment ID", nil))
	}

	deployment, err = u.deploymentRepository.FindOne(ctx, deploymentID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find deployment by ID", nil))
		}
		return nil, err // Return the NotFoundError directly
	}

	return deployment, nil
}
