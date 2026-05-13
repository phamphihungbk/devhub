package deploymentrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Deployment struct {
	model.Deployments
}

func (d *Deployment) ToEntity() *entity.Deployment {
	status, err := new(entity.DeploymentStatus).Parse(d.Status)
	if err != nil {
		return nil
	}

	return &entity.Deployment{
		ID:            d.ID,
		ServiceID:     d.ServiceID,
		EnvironmentID: d.EnvironmentID,
		ReleaseID:     d.ReleaseID,
		PluginID:      d.PluginID,
		Status:        status,
		ExternalRef:   misc.GetValue(d.ExternalRef),
		CommitSHA:     misc.GetValue(d.CommitSha),
		TriggeredBy:   d.TriggeredBy,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
		StartedAt:     misc.DerefTime(d.StartedAt),
		FinishedAt:    misc.DerefTime(d.FinishedAt),
	}
}

type Deployments []Deployment

func (ds Deployments) ToEntities() *entity.Deployments {
	deployments := make(entity.Deployments, 0, len(ds))
	for _, c := range ds {
		deployment := c.ToEntity()
		if deployment == nil {
			continue
		}
		deployments = append(deployments, misc.GetValue(deployment))
	}

	return misc.ToPointer(deployments)
}
