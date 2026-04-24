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
	FindAllApprovalRequests(ctx context.Context, filter FindAllApprovalRequestsFilter) (*entity.ApprovalRequests, int64, error)
	FindApprovalRequest(ctx context.Context, id uuid.UUID) (*entity.ApprovalRequest, error)
	FindApprovalDecisions(ctx context.Context, approvalRequestID uuid.UUID) ([]entity.ApprovalDecision, error)
	UpdateApprovalRequest(ctx context.Context, input UpdateApprovalRequestInput) (*entity.ApprovalRequest, error)
	CreateApprovalDecision(ctx context.Context, decision *entity.ApprovalDecision) (*entity.ApprovalDecision, error)
}

type FindAllApprovalRequestsFilter struct {
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
	Status    *entity.ApprovalRequestStatus
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
