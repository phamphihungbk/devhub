package service

import (
	"encoding/json"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"

	"github.com/google/uuid"
)

type ServiceDependency struct {
	model.ServiceDependencies
}

func (d *ServiceDependency) ToEntity() *entity.ServiceDependency {
	if d == nil {
		return nil
	}

	config := map[string]any{}
	if len(d.Config) > 0 {
		_ = json.Unmarshal([]byte(d.Config), &config)
	}

	dependency := &entity.ServiceDependency{
		ID:                 d.ID,
		ServiceID:          d.ServiceID,
		DependsOnServiceID: d.DependsOnServiceID,
		Type:               d.Type,
		Protocol:           misc.GetValue(d.Protocol),
		Port:               dependencyPortToEntity(d.Port),
		Path:               misc.GetValue(d.Path),
		Config:             config,
		CreatedBy:          d.CreatedBy,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}

	return dependency
}

type ServiceDependencies []ServiceDependency

func (ds ServiceDependencies) ToEntities() *entity.ServiceDependencies {
	dependencies := make(entity.ServiceDependencies, 0, len(ds))
	for _, d := range ds {
		dependency := d.ToEntity()
		if dependency == nil {
			continue
		}
		dependencies = append(dependencies, misc.GetValue(dependency))
	}

	return misc.ToPointer(dependencies)
}

type ServiceDependencyWithService struct {
	ServiceDependency
	DependsOnProjectID *uuid.UUID `db:"depends_on_project_id"`
	DependsOnName      *string    `db:"depends_on_name"`
	DependsOnRepoURL   *string    `db:"depends_on_repo_url"`
}

func (d *ServiceDependencyWithService) ToEntity() *entity.ServiceDependency {
	if d == nil {
		return nil
	}

	dependency := d.ServiceDependency.ToEntity()
	if dependency == nil {
		return nil
	}

	if d.DependsOnProjectID != nil && d.DependsOnName != nil && d.DependsOnRepoURL != nil {
		dependency.DependsOnService = &entity.Service{
			ID:        d.DependsOnServiceID,
			ProjectID: *d.DependsOnProjectID,
			Name:      *d.DependsOnName,
			RepoURL:   *d.DependsOnRepoURL,
		}
	}

	return dependency
}

type ServiceDependenciesWithService []ServiceDependencyWithService

func (ds ServiceDependenciesWithService) ToEntities() *entity.ServiceDependencies {
	dependencies := make(entity.ServiceDependencies, 0, len(ds))
	for _, d := range ds {
		dependency := d.ToEntity()
		if dependency == nil {
			continue
		}
		dependencies = append(dependencies, misc.GetValue(dependency))
	}

	return misc.ToPointer(dependencies)
}

func dependencyPortToEntity(port *int32) *int {
	if port == nil {
		return nil
	}

	value := int(*port)
	return &value
}
