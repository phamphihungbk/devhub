package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type ApprovalRepository interface {
	CreateApprovalPolicy(ctx context.Context, policy *entity.ApprovalPolicy) (*entity.ApprovalPolicy, error)
	FindMatchingApprovalPolicy(ctx context.Context, input FindMatchingApprovalPolicyInput) (*entity.ApprovalPolicy, error)
	CreateApprovalRequest(ctx context.Context, request *entity.ApprovalRequest) (*entity.ApprovalRequest, error)
	FindApprovalRequest(ctx context.Context, id uuid.UUID) (*entity.ApprovalRequest, error)
	UpdateApprovalRequest(ctx context.Context, input UpdateApprovalRequestInput) (*entity.ApprovalRequest, error)
	CreateApprovalDecision(ctx context.Context, decision *entity.ApprovalDecision) (*entity.ApprovalDecision, error)
}

type FindMatchingApprovalPolicyInput struct {
	Resource    string
	Action      string
	ProjectID   *uuid.UUID
	ServiceID   *uuid.UUID
	Environment *string
}

type UpdateApprovalRequestInput struct {
	ID                uuid.UUID
	Status            *entity.ApprovalRequestStatus
	ApprovedCount     *int
	RejectedCount     *int
	RequiredApprovals *int
	ResolvedAt        *time.Time
}
