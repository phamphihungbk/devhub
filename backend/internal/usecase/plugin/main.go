package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type PluginUsecase interface {
	CreatePlugin(ctx context.Context, user CreatePluginInput) (*entity.Plugin, error)
	FindOnePlugin(ctx context.Context, id FindOnePluginInput) (*entity.Plugin, error)
	FindAllPlugins(ctx context.Context, input FindAllPluginsInput) (entity.Page[entity.Plugin], error)
	UpdatePlugin(ctx context.Context, input UpdatePluginInput) (*entity.Plugin, error)
	DeletePlugin(ctx context.Context, id DeletePluginInput) (*entity.Plugin, error)
}

type pluginUsecase struct {
	appConfig        config.AppConfig
	pluginRepository repository.PluginRepository
}

func NewPluginUsecase(
	appConfig config.AppConfig,
	pluginRepository repository.PluginRepository,
) PluginUsecase {
	return &pluginUsecase{
		appConfig:        appConfig,
		pluginRepository: pluginRepository,
	}
}
