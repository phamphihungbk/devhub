package usecase

import (
	"context"

	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ServiceUsecase interface {
	FindAllServices(ctx context.Context, input FindAllServicesInput) (entity.Services, error)
	FindServiceDependencies(ctx context.Context, input FindServiceDependenciesInput) (entity.ServiceDependencies, error)
	CreateServiceDependency(ctx context.Context, input CreateServiceDependencyInput) (*entity.ServiceDependency, error)
	DeleteServiceDependency(ctx context.Context, input DeleteServiceDependencyInput) (*entity.ServiceDependency, error)
}

type serviceUsecase struct {
	appConfig         config.AppConfig
	serviceRepository repository.ServiceRepository
}

func NewServiceUsecase(
	appConfig config.AppConfig,
	serviceRepository repository.ServiceRepository,
) ServiceUsecase {
	return &serviceUsecase{
		appConfig:         appConfig,
		serviceRepository: serviceRepository,
	}
}
