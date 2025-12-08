package usecase

import (
	"context"
	"devhub-backend/internal/domain/entity"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/misc"

	"devhub-backend/pkg/validator"
)

type CreatePluginInput struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin user"`
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

	plugin = &entity.Plugin{
		Name:  input.Name,
		Email: input.Email,
		Role:  entity.PluginRole(input.Role),
	}

	created, err := u.pluginRepository.CreateOne(ctx, plugin)
	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to create plugin", nil))
	}

	return created, nil
}
