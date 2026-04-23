package usecase

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type DeploymentUsecase interface {
	CreateDeployment(ctx context.Context, deployment CreateDeploymentInput) (*entity.Deployment, error)
	FindOneDeployment(ctx context.Context, id FindOneDeploymentInput) (*entity.Deployment, error)
	FindAllDeployments(ctx context.Context, input FindAllDeploymentsInput) (entity.Page[entity.Deployment], error)
	UpdateDeployment(ctx context.Context, input UpdateDeploymentInput) (*entity.Deployment, error)
	DeleteDeployment(ctx context.Context, id DeleteDeploymentInput) (*entity.Deployment, error)
}

type deploymentUsecase struct {
	approvalRepository   repository.ApprovalRepository
	appConfig            config.AppConfig
	deploymentRepository repository.DeploymentRepository
}

func NewDeploymentUsecase(
	appConfig config.AppConfig,
	approvalRepository repository.ApprovalRepository,
	deploymentRepository repository.DeploymentRepository,
) DeploymentUsecase {
	return &deploymentUsecase{
		approvalRepository:   approvalRepository,
		appConfig:            appConfig,
		deploymentRepository: deploymentRepository,
	}
}
