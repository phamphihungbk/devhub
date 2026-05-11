package jobrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Job struct {
	model.Jobs
}

func (j *Job) ToEntity() *entity.Job {
	jobType, err := new(entity.JobType).Parse(j.Type)
	if err != nil {
		return nil
	}
	status, err := new(entity.JobStatus).Parse(j.Status)
	if err != nil {
		return nil
	}
	resourceType, err := new(entity.JobResourceType).Parse(j.ResourceType)
	if err != nil {
		return nil
	}
	payload, err := new(entity.JobPayload).Parse(j.Payload)
	if err != nil {
		return nil
	}

	return &entity.Job{
		ID:           j.ID,
		Type:         jobType,
		Status:       status,
		ResourceType: resourceType,
		ResourceID:   j.ResourceID,
		PluginID:     j.PluginID,
		Payload:      payload,
		Result:       j.Result,
		Error:        misc.GetValue(j.Error),
		Attempts:     int(j.Attempts),
		MaxAttempts:  int(j.MaxAttempts),
		CreatedBy:    j.CreatedBy,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
		StartedAt:    misc.DerefTime(j.StartedAt),
		FinishedAt:   misc.DerefTime(j.FinishedAt),
	}
}
