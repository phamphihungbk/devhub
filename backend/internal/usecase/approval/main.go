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
}

type approvalUsecase struct {
	appConfig                 config.AppConfig
	approvalRepository        repository.ApprovalRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

func NewApprovalUsecase(
	appConfig config.AppConfig,
	approvalRepository repository.ApprovalRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
) ApprovalUsecase {
	return &approvalUsecase{
		appConfig:                 appConfig,
		approvalRepository:        approvalRepository,
		scaffoldRequestRepository: scaffoldRequestRepository,
	}
}
