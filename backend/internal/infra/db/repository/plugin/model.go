package userrepo

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

	return &entity.Plugin{
		ID:          c.ID,
		Type:        pluginType,
		Name:        c.Name,
		Version:     c.Version,
		Description: misc.GetValue(c.Description),
		InstalledAt: c.InstalledAt,
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
