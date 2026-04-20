package approvalrepo

import (
	"devhub-backend/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type approvalPolicyModel struct {
	ID                uuid.UUID  `db:"id"`
	Resource          string     `db:"resource"`
	Action            string     `db:"action"`
	ProjectID         *uuid.UUID `db:"project_id"`
	ServiceID         *uuid.UUID `db:"service_id"`
	Environment       *string    `db:"environment"`
	RequiredApprovals int        `db:"required_approvals"`
	Enabled           bool       `db:"enabled"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

func (m approvalPolicyModel) toEntity() *entity.ApprovalPolicy {
	return &entity.ApprovalPolicy{
		ID:                m.ID,
		Resource:          m.Resource,
		Action:            m.Action,
		ProjectID:         m.ProjectID,
		ServiceID:         m.ServiceID,
		Environment:       m.Environment,
		RequiredApprovals: m.RequiredApprovals,
		Enabled:           m.Enabled,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

type approvalRequestModel struct {
	ID                uuid.UUID  `db:"id"`
	Resource          string     `db:"resource"`
	Action            string     `db:"action"`
	ResourceID        uuid.UUID  `db:"resource_id"`
	RequestedBy       uuid.UUID  `db:"requested_by"`
	ProjectID         *uuid.UUID `db:"project_id"`
	ServiceID         *uuid.UUID `db:"service_id"`
	Environment       *string    `db:"environment"`
	Status            string     `db:"status"`
	RequiredApprovals int        `db:"required_approvals"`
	ApprovedCount     int        `db:"approved_count"`
	RejectedCount     int        `db:"rejected_count"`
	ResolvedAt        *time.Time `db:"resolved_at"`
}

func (m approvalRequestModel) toEntity() *entity.ApprovalRequest {
	status, err := new(entity.ApprovalRequestStatus).Parse(m.Status)
	if err != nil {
		return nil
	}

	return &entity.ApprovalRequest{
		ID:                m.ID,
		Resource:          m.Resource,
		Action:            m.Action,
		ResourceID:        m.ResourceID,
		RequestedBy:       m.RequestedBy,
		ProjectID:         m.ProjectID,
		ServiceID:         m.ServiceID,
		Environment:       m.Environment,
		Status:            status,
		RequiredApprovals: m.RequiredApprovals,
		ApprovedCount:     m.ApprovedCount,
		RejectedCount:     m.RejectedCount,
		ResolvedAt:        m.ResolvedAt,
	}
}

type approvalDecisionModel struct {
	ID                uuid.UUID `db:"id"`
	ApprovalRequestID uuid.UUID `db:"approval_request_id"`
	DecidedBy         uuid.UUID `db:"decided_by"`
	Decision          string    `db:"decision"`
	Comment           string    `db:"comment"`
	CreatedAt         time.Time `db:"created_at"`
}

func (m approvalDecisionModel) toEntity() *entity.ApprovalDecision {
	decision, err := new(entity.ApprovalDecisionType).Parse(m.Decision)
	if err != nil {
		return nil
	}

	return &entity.ApprovalDecision{
		ID:                m.ID,
		ApprovalRequestID: m.ApprovalRequestID,
		DecidedBy:         m.DecidedBy,
		Decision:          decision,
		Comment:           m.Comment,
		CreatedAt:         m.CreatedAt,
	}
}
