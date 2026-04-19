package usecase

import (
	"context"

	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ServiceUsecase interface {
	FindAllServices(ctx context.Context, input FindAllServicesInput) (entity.Services, error)
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
