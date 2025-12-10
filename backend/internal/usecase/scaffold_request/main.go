package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ScaffoldRequestUsecase interface {
	CreateScaffoldRequest(ctx context.Context, input CreateScaffoldRequestInput) (*entity.ScaffoldRequest, error)
	FindOneScaffoldRequest(ctx context.Context, id FindOneScaffoldRequestInput) (*entity.ScaffoldRequest, error)
	FindAllScaffoldRequests(ctx context.Context, input FindAllScaffoldRequestsInput) (entity.Page[entity.ScaffoldRequest], error)
	DeleteScaffoldRequest(ctx context.Context, id DeleteScaffoldRequestInput) (*entity.ScaffoldRequest, error)
}

type scaffoldRequestUsecase struct {
	appConfig                 config.AppConfig
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

func NewScaffoldRequestUsecase(
	appConfig config.AppConfig,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
) ScaffoldRequestUsecase {
	return &scaffoldRequestUsecase{
		appConfig:                 appConfig,
		scaffoldRequestRepository: scaffoldRequestRepository,
	}
}
