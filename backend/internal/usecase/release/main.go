package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	releaseprovider "devhub-backend/internal/infra/release_provider"
)

type ReleaseUsecase interface {
	CreateRelease(ctx context.Context, input CreateReleaseInput) (*entity.Release, error)
}

type releaseUsecase struct {
	appConfig         config.AppConfig
	projectRepository repository.ProjectRepository
	releaseRepository repository.ReleaseRepository
	releaseClients    map[string]releaseprovider.Client
}

func NewReleaseUsecase(
	appConfig config.AppConfig,
	projectRepository repository.ProjectRepository,
	releaseRepository repository.ReleaseRepository,
	giteaCfg config.GiteaConfig,
) ReleaseUsecase {
	giteaClient := releaseprovider.NewClient(giteaCfg)

	return &releaseUsecase{
		appConfig:         appConfig,
		projectRepository: projectRepository,
		releaseRepository: releaseRepository,
		releaseClients: map[string]releaseprovider.Client{
			giteaClient.Provider(): giteaClient,
		},
	}
}
