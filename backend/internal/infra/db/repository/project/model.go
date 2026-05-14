package projectrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Project struct {
	model.Projects
}

func (c *Project) ToEntity() *entity.Project {
	return &entity.Project{
		ID:          c.ID,
		Name:        c.Name,
		Description: misc.GetValue(c.Description),
		OwnerTeamID: c.OwnerTeamID,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type Projects []Project

func (ps Projects) ToEntities() *entity.Projects {
	projects := make(entity.Projects, 0, len(ps))
	for _, c := range ps {
		project := c.ToEntity()
		if project == nil {
			continue
		}
		projects = append(projects, misc.GetValue(project))
	}

	return misc.ToPointer(projects)
}
