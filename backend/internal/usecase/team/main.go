package usecase

import (
	"context"

	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type TeamUsecase interface {
	CreateTeam(ctx context.Context, input CreateTeamInput) (*entity.Team, error)
	UpdateTeam(ctx context.Context, input UpdateTeamInput) (*entity.Team, error)
	FindAllTeams(ctx context.Context, input FindAllTeamsInput) (entity.Page[entity.Team], error)
}

type teamUsecase struct {
	appConfig      config.AppConfig
	teamRepository repository.TeamRepository
}

func NewTeamUsecase(appConfig config.AppConfig, teamRepository repository.TeamRepository) TeamUsecase {
	return &teamUsecase{
		appConfig:      appConfig,
		teamRepository: teamRepository,
	}
}
