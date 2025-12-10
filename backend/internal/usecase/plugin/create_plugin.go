package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"
)

type CreatePluginInput struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Version     string `json:"version" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=scaffolder runner"`
	Description string `json:"description" validate:"required,min=2,max=100"`
}

func (u *pluginUsecase) CreatePlugin(ctx context.Context, input CreatePluginInput) (plugin *entity.Plugin, err error) {
	const errLocation = "[usecase plugin/create_plugin CreatePlugin] "
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

	pluginType, err := new(entity.PluginType).Parse(input.Type)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewBadRequestError("invalid plugin type", nil))
	}

	plugin = &entity.Plugin{
		Name:        input.Name,
		Version:     input.Version,
		Type:        pluginType,
		Description: input.Description,
	}

	created, err := u.pluginRepository.CreateOne(ctx, plugin)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create plugin", nil))
	}

	return created, nil
}
