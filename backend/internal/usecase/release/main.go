package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ReleaseUsecase interface {
	CreateRelease(ctx context.Context, input CreateReleaseInput) (*entity.Release, error)
	FindAllReleases(ctx context.Context, input FindAllReleasesInput) (entity.Releases, error)
}

type releaseUsecase struct {
	pluginRepository  repository.PluginRepository
	appConfig         config.AppConfig
	releaseRepository repository.ReleaseRepository
}

func NewReleaseUsecase(
	appConfig config.AppConfig,
	pluginRepository repository.PluginRepository,
	releaseRepository repository.ReleaseRepository,
) ReleaseUsecase {

	return &releaseUsecase{
		pluginRepository:  pluginRepository,
		appConfig:         appConfig,
		releaseRepository: releaseRepository,
	}
}
