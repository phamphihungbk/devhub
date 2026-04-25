package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ScaffoldRequestUsecase interface {
	CreateScaffoldRequest(ctx context.Context, input CreateScaffoldRequestInput) (*entity.ScaffoldRequest, error)
	SuggestScaffoldRequest(ctx context.Context, input SuggestScaffoldRequestInput) (ScaffoldRequestSuggestion, error)
	FindOneScaffoldRequest(ctx context.Context, id FindOneScaffoldRequestInput) (*entity.ScaffoldRequest, error)
	FindAllScaffoldRequests(ctx context.Context, input FindAllScaffoldRequestsInput) (entity.Page[entity.ScaffoldRequest], error)
	DeleteScaffoldRequest(ctx context.Context, id DeleteScaffoldRequestInput) (*entity.ScaffoldRequest, error)
}

type scaffoldRequestUsecase struct {
	appConfig                 config.AppConfig
	approvalRepository        repository.ApprovalRepository
	pluginRepository          repository.PluginRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

func NewScaffoldRequestUsecase(
	appConfig config.AppConfig,
	approvalRepository repository.ApprovalRepository,
	pluginRepository repository.PluginRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
) ScaffoldRequestUsecase {
	return &scaffoldRequestUsecase{
		approvalRepository:        approvalRepository,
		appConfig:                 appConfig,
		pluginRepository:          pluginRepository,
		scaffoldRequestRepository: scaffoldRequestRepository,
	}
}
