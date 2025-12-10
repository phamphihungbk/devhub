package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"

	"github.com/google/uuid"
)

type DeletePluginInput struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (u *pluginUsecase) DeletePlugin(ctx context.Context, input DeletePluginInput) (plugin *entity.Plugin, err error) {
	const errLocation = "[usecase plugin/delete_plugin DeletePlugin] "
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

	deleted, err := u.pluginRepository.DeleteOne(ctx, pluginID)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to delete plugin", nil))
	}

	return deleted, nil
}
