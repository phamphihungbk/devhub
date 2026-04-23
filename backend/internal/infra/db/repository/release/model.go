package releaserepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Release struct {
	model.Releases
}

func (r *Release) ToEntity() *entity.Release {
	return &entity.Release{
		ID:          r.ID,
		ServiceID:   r.ServiceID,
		PluginID:    r.PluginID,
		Tag:         r.Tag,
		Target:      r.Target,
		Name:        r.Name,
		Status:      entity.ReleaseStatus(r.Status),
		Notes:       r.Notes,
		HTMLURL:     r.HTMLURL,
		ExternalRef: r.ExternalRef,
		TriggeredBy: r.TriggeredBy,
		CreatedAt:   r.CreatedAt,
	}
}

type Releases []Release

func (rs Releases) ToEntities() *entity.Releases {
	releases := make(entity.Releases, 0, len(rs))
	for _, r := range rs {
		release := r.ToEntity()
		if release == nil {
			continue
		}
		releases = append(releases, misc.GetValue(release))
	}

	return misc.ToPointer(releases)
}
