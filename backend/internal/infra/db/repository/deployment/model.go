package userrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Deployment struct {
	model.Deployments
}

func (c *Deployment) ToEntity() *entity.Deployment {
	return &entity.Deployment{
		ID:        c.ID,
		Email:     c.Email,
		Name:      c.Name,
		Role:      entity.UserRole(c.Role),
		CreatedAt: c.CreatedAt,
		LastLogin: misc.DerefTime(c.LastLogin),
		DeletedAt: misc.DerefTime(c.DeletedAt),
	}
}

type Deployments []Deployment

func (us Deployments) ToEntities() *entity.Deployments {
	deployments := make(entity.Deployments, 0, len(us))
	for _, c := range us {
		deployment := c.ToEntity()
		if deployment == nil {
			continue
		}
		deployments = append(deployments, misc.GetValue(deployment))
	}

	return misc.ToPointer(deployments)
}
