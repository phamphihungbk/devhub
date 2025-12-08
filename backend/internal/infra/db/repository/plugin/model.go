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
	return &entity.Plugin{
		ID:        c.ID,
		Email:     c.Email,
		Name:      c.Name,
		Role:      entity.UserRole(c.Role),
		CreatedAt: c.CreatedAt,
		LastLogin: misc.DerefTime(c.LastLogin),
		DeletedAt: misc.DerefTime(c.DeletedAt),
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
