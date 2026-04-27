package usecase

import (
	"context"

	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
)

type ApprovalUsecase interface {
	CreateApprovalPolicy(ctx context.Context, input CreateApprovalPolicyInput) (*entity.ApprovalPolicy, error)
	CreateApprovalDecision(ctx context.Context, input CreateApprovalDecisionInput) (*CreateApprovalDecisionOutput, error)
	FindAllApprovalRequests(ctx context.Context, input FindAllApprovalRequestsInput) (entity.Page[entity.ApprovalRequest], error)
	FindApprovalRequestDetail(ctx context.Context, input FindApprovalRequestDetailInput) (*FindApprovalRequestDetailOutput, error)
}

type approvalUsecase struct {
	appConfig                 config.AppConfig
	approvalRepository        repository.ApprovalRepository
	deploymentRepository      repository.DeploymentRepository
	projectRepository         repository.ProjectRepository
	releaseRepository         repository.ReleaseRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
	serviceRepository         repository.ServiceRepository
	userRepository            repository.UserRepository
}

func NewApprovalUsecase(
	appConfig config.AppConfig,
	approvalRepository repository.ApprovalRepository,
	deploymentRepository repository.DeploymentRepository,
	projectRepository repository.ProjectRepository,
	releaseRepository repository.ReleaseRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
	serviceRepository repository.ServiceRepository,
	userRepository repository.UserRepository,
) ApprovalUsecase {
	return &approvalUsecase{
		appConfig:                 appConfig,
		approvalRepository:        approvalRepository,
		deploymentRepository:      deploymentRepository,
		projectRepository:         projectRepository,
		releaseRepository:         releaseRepository,
		scaffoldRequestRepository: scaffoldRequestRepository,
		serviceRepository:         serviceRepository,
		userRepository:            userRepository,
	}
}
