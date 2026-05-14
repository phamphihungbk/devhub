package validator

import (
	"context"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sync"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PluginConfigSchemaValidator struct {
	Context          context.Context
	PluginRepository repository.PluginRepository
	mu               sync.Mutex
	message          string
}

func (*PluginConfigSchemaValidator) Tag() string {
	return "plugin_config_schema"
}

func (v *PluginConfigSchemaValidator) Func() validator.Func {
	return func(fl validator.FieldLevel) bool {
		if v.PluginRepository == nil {
			v.setMessage("plugin repository is not configured")
			return false
		}

		pluginID, ok := pluginIDFromTop(fl.Top())
		if !ok {
			return true
		}

		ctx := v.Context
		if ctx == nil {
			ctx = context.Background()
		}

		plugin, err := v.PluginRepository.FindOne(ctx, pluginID)
		if err != nil {
			return true
		}
		if plugin.ConfigSchema.IsZero() {
			return true
		}

		variables, ok := fl.Field().Interface().(map[string]interface{})
		if !ok {
			v.setMessage("variables must be an object")
			return false
		}

		if message := validateObjectSchema(variables, plugin.ConfigSchema.Required, plugin.ConfigSchema.Properties, plugin.ConfigSchema.AdditionalProperties, "variables"); message != "" {
			v.setMessage(message)
			return false
		}

		v.setMessage("")
		return true
	}
}

func (v *PluginConfigSchemaValidator) Translation() (string, validator.TranslationFunc) {
	return "{0}",
		func(ut ut.Translator, fe validator.FieldError) string {
			detail := v.getMessage()
			if detail == "" {
				detail = fe.Field() + " must match the selected plugin config schema"
			}
			return detail
		}
}

func (v *PluginConfigSchemaValidator) setMessage(message string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.message = message
}

func (v *PluginConfigSchemaValidator) getMessage() string {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.message
}

func pluginIDFromTop(top reflect.Value) (uuid.UUID, bool) {
	if top.Kind() == reflect.Pointer {
		top = top.Elem()
	}
	if top.Kind() != reflect.Struct {
		return uuid.UUID{}, false
	}

	field := top.FieldByName("PluginID")
	if !field.IsValid() || field.Kind() != reflect.String {
		return uuid.UUID{}, false
	}

	pluginID, err := uuid.Parse(field.String())
	return pluginID, err == nil
}

func validateObjectSchema(value map[string]interface{}, required []string, properties map[string]entity.PluginSchemaProperty, additionalProperties *bool, path string) string {
	for _, key := range required {
		if _, ok := value[key]; !ok {
			return fmt.Sprintf("%s.%s is required", path, key)
		}
	}

	for key, raw := range value {
		fieldPath := fmt.Sprintf("%s.%s", path, key)
		property, ok := properties[key]
		if !ok {
			if additionalProperties == nil || *additionalProperties {
				continue
			}
			return fmt.Sprintf("%s is not allowed", fieldPath)
		}
		if raw == nil {
			return fmt.Sprintf("%s must not be null", fieldPath)
		}
		if message := validateSchemaProperty(raw, property, fieldPath); message != "" {
			return message
		}
	}

	return ""
}

func validateSchemaProperty(value interface{}, property entity.PluginSchemaProperty, path string) string {
	if len(property.Enum) > 0 && !matchesEnum(value, property.Enum) {
		return fmt.Sprintf("%s must be one of: %v", path, property.Enum)
	}

	switch property.Type {
	case "", "any":
		return ""
	case "string":
		text, ok := value.(string)
		if !ok {
			return fmt.Sprintf("%s must be a string", path)
		}
		if property.MinLength != nil && len(text) < *property.MinLength {
			return fmt.Sprintf("%s must be at least %d characters", path, *property.MinLength)
		}
		if property.Pattern != "" {
			matched, err := regexp.MatchString(property.Pattern, text)
			if err != nil || !matched {
				return fmt.Sprintf("%s must match pattern %q", path, property.Pattern)
			}
		}
		return ""
	case "integer":
		number, ok := numberValue(value)
		if !ok || math.Trunc(number) != number {
			return fmt.Sprintf("%s must be an integer", path)
		}
		return numberBoundsMessage(number, property, path)
	case "number":
		number, ok := numberValue(value)
		if !ok {
			return fmt.Sprintf("%s must be a number", path)
		}
		return numberBoundsMessage(number, property, path)
	case "boolean":
		_, ok := value.(bool)
		if !ok {
			return fmt.Sprintf("%s must be a boolean", path)
		}
		return ""
	case "object":
		nested, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Sprintf("%s must be an object", path)
		}
		return validateObjectSchema(nested, property.Required, property.Properties, property.AdditionalProperties, path)
	case "array":
		_, ok := value.([]interface{})
		if !ok {
			return fmt.Sprintf("%s must be an array", path)
		}
		return ""
	default:
		return fmt.Sprintf("%s has unsupported schema type %q", path, property.Type)
	}
}

func matchesEnum(value interface{}, enum []string) bool {
	for _, allowed := range enum {
		if fmt.Sprint(value) == allowed {
			return true
		}
	}
	return false
}

func numberValue(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func numberBoundsMessage(value float64, property entity.PluginSchemaProperty, path string) string {
	if property.Minimum != nil && value < *property.Minimum {
		return fmt.Sprintf("%s must be greater than or equal to %g", path, *property.Minimum)
	}
	if property.Maximum != nil && value > *property.Maximum {
		return fmt.Sprintf("%s must be less than or equal to %g", path, *property.Maximum)
	}
	return ""
}
