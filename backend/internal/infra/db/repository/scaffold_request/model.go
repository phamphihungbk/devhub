package projectrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type ScaffoldRequest struct {
	model.ScaffoldRequests
}

func (c *ScaffoldRequest) ToEntity() *entity.ScaffoldRequest {
	env, err := new(entity.ProjectEnvironment).Parse(c.Environment)
	if err != nil {
		return nil
	}
	variables, err := new(entity.ScaffoldRequestVariables).Parse(c.Variables)
	if err != nil {
		return nil
	}

	return &entity.ScaffoldRequest{
		ID:          c.ID,
		ProjectID:   c.ProjectID,
		Template:    c.Template,
		Environment: env,
		Variables:   variables,
	}
}

type ScaffoldRequests []ScaffoldRequest

func (ps ScaffoldRequests) ToEntities() *entity.ScaffoldRequests {
	scaffoldRequests := make(entity.ScaffoldRequests, 0, len(ps))
	for _, c := range ps {
		scaffoldRequest := c.ToEntity()
		if scaffoldRequest == nil {
			continue
		}
		scaffoldRequests = append(scaffoldRequests, misc.GetValue(scaffoldRequest))
	}

	return misc.ToPointer(scaffoldRequests)
}
