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

type FindOnePluginInput struct {
	ID string `json:"id" validate:"required,uuid4"`
}

func (u *pluginUsecase) FindOnePlugin(ctx context.Context, input FindOnePluginInput) (plugin *entity.Plugin, err error) {
	const errLocation = "[usecase plugin/find_one_plugin FindOnePlugin] "
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

	pluginID, err := uuid.Parse(input.ID)

	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid plugin ID", nil))
	}

	plugin, err = u.pluginRepository.FindOne(ctx, pluginID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return nil, misc.WrapError(err, errs.NewInternalServerError("failed to find plugin by ID", nil))
		}
		return nil, err // Return the NotFoundError directly
	}

	return plugin, nil
}
