package service

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Service struct {
	model.Services
}

func (c *Service) ToEntity() *entity.Service {
	return &entity.Service{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		Name:      c.Name,
		RepoURL:   c.RepoURL,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

type Services []Service

func (ps Services) ToEntities() *entity.Services {
	services := make(entity.Services, 0, len(ps))
	for _, c := range ps {
		service := c.ToEntity()
		if service == nil {
			continue
		}
		services = append(services, misc.GetValue(service))
	}

	return misc.ToPointer(services)
}
