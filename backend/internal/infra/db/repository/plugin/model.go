package pluginrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Plugin struct {
	model.Plugins
}

func (c *Plugin) ToEntity() *entity.Plugin {
	pluginType, err := new(entity.PluginType).Parse(c.Type)
	if err != nil {
		return nil
	}
	pluginRuntime, err := new(entity.PluginRuntime).Parse(c.Runtime)
	if err != nil {
		return nil
	}

	var configSchema entity.PluginConfigSchema
	if c.ConfigSchema != nil {
		configSchema, err = new(entity.PluginConfigSchema).Parse(*c.ConfigSchema)
		if err != nil {
			return nil
		}
	}

	return &entity.Plugin{
		ID:           c.ID,
		Name:         c.Name,
		Version:      c.Version,
		Type:         pluginType,
		Runtime:      pluginRuntime,
		Entrypoint:   c.Entrypoint,
		ConfigSchema: configSchema,
		Description:  misc.GetValue(c.Description),
		Enabled:      c.Enabled,
		CreatedAt:    c.CreatedAt,
	}
}

type Plugins []Plugin

func (us Plugins) ToEntities() *entity.Plugins {
	plugins := make(entity.Plugins, 0, len(us))
	for _, c := range us {
		plugin := c.ToEntity()
		if plugin == nil {
			continue
		}
		plugins = append(plugins, misc.GetValue(plugin))
	}

	return misc.ToPointer(plugins)
}
