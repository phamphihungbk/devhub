package validator

import (
	"context"
	"devhub-backend/internal/domain/repository"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EnabledPluginValidator struct {
	Context          context.Context
	PluginRepository repository.PluginRepository
}

func (EnabledPluginValidator) Tag() string {
	return "enabled_plugin"
}

func (v EnabledPluginValidator) Func() validator.Func {
	return func(fl validator.FieldLevel) bool {
		if v.PluginRepository == nil {
			return false
		}

		ctx := v.Context
		if ctx == nil {
			ctx = context.Background()
		}

		pluginID, err := uuid.Parse(fl.Field().String())
		if err != nil {
			return false
		}

		plugin, err := v.PluginRepository.FindOne(ctx, pluginID)
		if err != nil {
			return false
		}

		return plugin.Enabled
	}
}

func (EnabledPluginValidator) Translation() (string, validator.TranslationFunc) {
	return "{0} must reference an enabled plugin",
		func(ut ut.Translator, fe validator.FieldError) string {
			msg, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.Field() + " must reference an enabled plugin"
			}
			return msg
		}
}
