package environment

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Environment struct {
	model.Environments
}

func (e *Environment) ToEntity() *entity.Environment {
	tier, err := new(entity.EnvironmentTier).Parse(e.Tier)
	if err != nil {
		return nil
	}

	config, err := new(entity.EnvironmentConfig).Parse(e.Config)
	if err != nil {
		return nil
	}

	return &entity.Environment{
		ID:             e.ID,
		ProjectID:      e.ProjectID,
		Name:           e.Name,
		Tier:           tier,
		Cluster:        e.Cluster,
		Namespace:      e.Namespace,
		ArgocdInstance: e.ArgocdInstance,
		Config:         config,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

type Environments []Environment

func (es Environments) ToEntities() *entity.Environments {
	environments := make(entity.Environments, 0, len(es))
	for _, item := range es {
		environment := item.ToEntity()
		if environment == nil {
			continue
		}
		environments = append(environments, misc.GetValue(environment))
	}

	return misc.ToPointer(environments)
}
