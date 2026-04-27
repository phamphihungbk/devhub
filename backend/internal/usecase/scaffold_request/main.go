package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/infra/ai"
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
	aiClient                  ai.Client
	pluginRepository          repository.PluginRepository
	projectRepository         repository.ProjectRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

func NewScaffoldRequestUsecase(
	appConfig config.AppConfig,
	approvalRepository repository.ApprovalRepository,
	aiClient ai.Client,
	pluginRepository repository.PluginRepository,
	projectRepository repository.ProjectRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
) ScaffoldRequestUsecase {
	return &scaffoldRequestUsecase{
		approvalRepository:        approvalRepository,
		aiClient:                  aiClient,
		appConfig:                 appConfig,
		pluginRepository:          pluginRepository,
		projectRepository:         projectRepository,
		scaffoldRequestRepository: scaffoldRequestRepository,
	}
}
