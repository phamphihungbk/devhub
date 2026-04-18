package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ReleaseUsecase interface {
	CreateRelease(ctx context.Context, input CreateReleaseInput) (*entity.Release, error)
}

type releaseUsecase struct {
	appConfig         config.AppConfig
	releaseRepository repository.ReleaseRepository
}

func NewReleaseUsecase(
	appConfig config.AppConfig,
	releaseRepository repository.ReleaseRepository,
) ReleaseUsecase {

	return &releaseUsecase{
		appConfig:         appConfig,
		releaseRepository: releaseRepository,
	}
}
