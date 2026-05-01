package approvalrepo

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/infra/db/model_gen/devhub/public/model"
	"devhub-backend/internal/util/misc"
)

type ApprovalPolicy struct {
	model.ApprovalPolicies
}

func (ap *ApprovalPolicy) toEntity() *entity.ApprovalPolicy {
	return &entity.ApprovalPolicy{
		ID:                ap.ID,
		Resource:          ap.Resource,
		Action:            ap.Action,
		ProjectID:         ap.ProjectID,
		ServiceID:         ap.ServiceID,
		EnvironmentID:     ap.EnvironmentID,
		RequiredApprovals: int(ap.RequiredApprovals),
		Enabled:           ap.Enabled,
		CreatedAt:         ap.CreatedAt,
		UpdatedAt:         ap.UpdatedAt,
	}
}

type ApprovalRequest struct {
	model.ApprovalRequests
}

func (ar *ApprovalRequest) toEntity() *entity.ApprovalRequest {
	status, err := new(entity.ApprovalRequestStatus).Parse(ar.Status)
	if err != nil {
		return nil
	}

	return &entity.ApprovalRequest{
		ID:                ar.ID,
		Resource:          ar.Resource,
		Action:            ar.Action,
		ResourceID:        ar.ResourceID,
		RequestedBy:       ar.RequestedBy,
		Status:            status,
		RequiredApprovals: int(ar.RequiredApprovals),
		ApprovedCount:     int(ar.ApprovedCount),
		RejectedCount:     int(ar.RejectedCount),
		ResolvedAt:        misc.DerefTime(ar.ResolvedAt),
		CreatedAt:         ar.CreatedAt,
		UpdatedAt:         ar.UpdatedAt,
	}
}

type ApprovalRequests []ApprovalRequest

func (ars ApprovalRequests) ToEntities() *entity.ApprovalRequests {
	requests := make(entity.ApprovalRequests, 0, len(ars))
	for _, item := range ars {
		request := item.toEntity()
		if request == nil {
			continue
		}
		requests = append(requests, *request)
	}

	return &requests
}

type ApprovalDecision struct {
	model.ApprovalDecisions
}

type ApprovalDecisions []ApprovalDecision

func (ad *ApprovalDecision) toEntity() *entity.ApprovalDecision {
	decision, err := new(entity.ApprovalDecisionType).Parse(ad.Decision)
	if err != nil {
		return nil
	}

	return &entity.ApprovalDecision{
		ID:                ad.ID,
		ApprovalRequestID: ad.ApprovalRequestID,
		DecidedBy:         ad.DecidedBy,
		Decision:          decision,
		Comment:           misc.GetValue(ad.Comment),
		CreatedAt:         ad.CreatedAt,
	}
}

func (ads ApprovalDecisions) ToEntities() []entity.ApprovalDecision {
	decisions := make([]entity.ApprovalDecision, 0, len(ads))
	for _, item := range ads {
		decision := item.toEntity()
		if decision == nil {
			continue
		}
		decisions = append(decisions, *decision)
	}

	return decisions
}
