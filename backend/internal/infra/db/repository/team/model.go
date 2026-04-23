package teamrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type Team struct {
	model.Teams
}

type Teams []Team

func (t *Team) ToEntity() *entity.Team {
	return &entity.Team{
		ID:           t.ID,
		Name:         t.Name,
		OwnerContact: t.OwnerContact,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		DeletedAt:    misc.DerefTime(t.DeletedAt),
	}
}

func (t Teams) ToEntities() *entity.Teams {
	teams := make(entity.Teams, 0, len(t))
	for _, team := range t {
		entityTeam := team.ToEntity()
		if entityTeam != nil {
			teams = append(teams, *entityTeam)
		}
	}
	return &teams
}
