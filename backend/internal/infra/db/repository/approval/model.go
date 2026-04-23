package approvalrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type ApprovalPolicy struct {
	model.ApprovalPolicies
}

func (m *ApprovalPolicy) toEntity() *entity.ApprovalPolicy {
	return &entity.ApprovalPolicy{
		ID:                m.ID,
		Resource:          m.Resource,
		Action:            m.Action,
		ProjectID:         m.ProjectID,
		ServiceID:         m.ServiceID,
		Environment:       m.Environment,
		RequiredApprovals: int(m.RequiredApprovals),
		Enabled:           m.Enabled,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

type ApprovalRequest struct {
	model.ApprovalRequests
}

type ApprovalRequests []ApprovalRequest

func (m ApprovalRequests) ToEntities() *entity.ApprovalRequests {
	requests := make(entity.ApprovalRequests, 0, len(m))
	for _, item := range m {
		request := item.toEntity()
		if request == nil {
			continue
		}
		requests = append(requests, *request)
	}

	return &requests
}

func (m *ApprovalRequest) toEntity() *entity.ApprovalRequest {
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
		RequiredApprovals: int(m.RequiredApprovals),
		ApprovedCount:     int(m.ApprovedCount),
		RejectedCount:     int(m.RejectedCount),
		ResolvedAt:        m.ResolvedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

type ApprovalDecision struct {
	model.ApprovalDecisions
}

func (m *ApprovalDecision) toEntity() *entity.ApprovalDecision {
	decision, err := new(entity.ApprovalDecisionType).Parse(m.Decision)
	if err != nil {
		return nil
	}

	return &entity.ApprovalDecision{
		ID:                m.ID,
		ApprovalRequestID: m.ApprovalRequestID,
		DecidedBy:         m.DecidedBy,
		Decision:          decision,
		Comment:           misc.GetValue(m.Comment),
		CreatedAt:         m.CreatedAt,
	}
}
