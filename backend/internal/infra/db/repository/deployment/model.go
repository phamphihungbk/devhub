package deploymentrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Deployment struct {
	model.Deployments
}

func (c *Deployment) ToEntity() *entity.Deployment {
	status, err := new(entity.DeploymentStatus).Parse(c.Status)
	if err != nil {
		return nil
	}
	env, err := new(entity.ProjectEnvironment).Parse(c.Environment)
	if err != nil {
		return nil
	}

	return &entity.Deployment{
		ID:          c.ID,
		ProjectID:   c.ProjectID,
		PluginID:    c.PluginID,
		Environment: env,
		Service:     c.Service,
		Version:     c.Version,
		Status:      status,
		ExternalRef: misc.GetValue(c.ExternalRef),
		CommitSHA:   misc.GetValue(c.CommitSha),
		TriggeredBy: c.TriggeredBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		FinishedAt:  c.FinishedAt,
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
