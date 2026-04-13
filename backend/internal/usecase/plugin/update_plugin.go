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

type UpdatePluginInput struct {
	ID          string  `json:"id" validate:"required,uuid"`
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Type        *string `json:"type" validate:"omitempty,oneof=scaffolder runner"`
	Version     *string `json:"version" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
	Entrypoint  *string `json:"entrypoint" validate:"omitempty,min=1,max=500"`
	Scope       *string `json:"scope" validate:"omitempty,oneof=global project environment"`
	Enabled     *bool   `json:"enabled"`
}

func (u *pluginUsecase) UpdatePlugin(ctx context.Context, input UpdatePluginInput) (plugin *entity.Plugin, err error) {
	const errLocation = "[usecase plugin/update_plugin UpdatePlugin] "
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

	updated, err := u.pluginRepository.UpdateOne(ctx, repository.UpdatePluginInput{
		ID:          uuid.MustParse(input.ID),
		Name:        input.Name,
		Type:        (*entity.PluginType)(input.Type),
		Version:     input.Version,
		Description: input.Description,
		Entrypoint:  input.Entrypoint,
		Scope:       input.Scope,
		Enabled:     input.Enabled,
	})

	if err != nil {
		return nil, misc.WrapError(err, errs.NewInternalServerError("failed to update plugin", nil))
	}

	return updated, nil
}
