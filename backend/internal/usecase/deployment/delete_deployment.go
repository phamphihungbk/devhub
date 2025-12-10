package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type DeleteDeploymentInput struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (u *deploymentUsecase) DeleteDeployment(ctx context.Context, input DeleteDeploymentInput) (deployment *entity.Deployment, err error) {
	const errLocation = "[usecase deployment/delete_deployment DeleteDeployment] "
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

	deleted, err := u.deploymentRepository.DeleteOne(ctx, deploymentID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete deployment", nil))
	}

	return deleted, nil
}
