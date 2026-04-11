package scaffoldrequestrepo

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
	status, err := new(entity.ScaffoldRequestStatus).Parse(c.Status)
	if err != nil {
		return nil
	}
	variables, err := new(entity.ScaffoldRequestVariables).Parse(c.Variables)
	if err != nil {
		return nil
	}

	return &entity.ScaffoldRequest{
		ID:            c.ID,
		PluginID:      c.PluginID,
		ProjectID:     c.ProjectID,
		RequestedBy:   c.RequestedBy,
		Template:      c.Template,
		Status:        status,
		Environment:   env,
		Variables:     variables,
		ApprovedBy:    c.ApprovedBy,
		ResultRepoURL: misc.GetValue(c.ResultRepoURL),
		ApprovedAt:    c.ApprovedAt,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
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
