package service

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Service struct {
	model.Services
}

func (s *Service) ToEntity() *entity.Service {
	return &entity.Service{
		ID:        s.ID,
		ProjectID: s.ProjectID,
		Name:      s.Name,
		RepoURL:   s.RepoURL,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

type Services []Service

func (ss Services) ToEntities() *entity.Services {
	services := make(entity.Services, 0, len(ss))
	for _, c := range ss {
		service := c.ToEntity()
		if service == nil {
			continue
		}
		services = append(services, misc.GetValue(service))
	}

	return misc.ToPointer(services)
}
