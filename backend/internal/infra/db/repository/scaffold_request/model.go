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
	return &entity.ScaffoldRequest{
		ID:           c.ID,
		Name:         c.Name,
		Description:  misc.GetValue(c.Description),
		Environments: misc.GetValue(c.Environments),
		CreatedAt:    misc.DerefTime(c.CreatedAt),
		UpdatedAt:    misc.DerefTime(c.UpdatedAt),
		DeletedAt:    misc.DerefTime(c.DeletedAt),
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
