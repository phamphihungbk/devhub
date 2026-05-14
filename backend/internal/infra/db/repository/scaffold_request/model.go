package scaffoldrequestrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
	"encoding/json"
)

type ScaffoldRequest struct {
	model.ScaffoldRequests
}

func (sr *ScaffoldRequest) ToEntity() *entity.ScaffoldRequest {
	status, err := new(entity.ScaffoldRequestStatus).Parse(sr.Status)
	if err != nil {
		return nil
	}
	variables := make(map[string]interface{})
	err = json.Unmarshal([]byte(sr.Variables), &variables)
	if err != nil {
		return nil
	}

	return &entity.ScaffoldRequest{
		ID:            sr.ID,
		ProjectID:     sr.ProjectID,
		PluginID:      sr.PluginID,
		RequestedBy:   sr.RequestedBy,
		ApprovedBy:    sr.ApprovedBy,
		Status:        status,
		Variables:     variables,
		ResultRepoURL: misc.GetValue(sr.ResultRepoURL),
		ApprovedAt:    misc.DerefTime(sr.ApprovedAt),
		CreatedAt:     sr.CreatedAt,
		UpdatedAt:     sr.UpdatedAt,
	}
}

type ScaffoldRequests []ScaffoldRequest

func (srs ScaffoldRequests) ToEntities() *entity.ScaffoldRequests {
	scaffoldRequests := make(entity.ScaffoldRequests, 0, len(srs))
	for _, c := range srs {
		scaffoldRequest := c.ToEntity()
		if scaffoldRequest == nil {
			continue
		}
		scaffoldRequests = append(scaffoldRequests, misc.GetValue(scaffoldRequest))
	}

	return misc.ToPointer(scaffoldRequests)
}
