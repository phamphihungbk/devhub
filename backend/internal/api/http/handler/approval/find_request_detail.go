package handler

import (
	"time"

	approvalUsecase "devhub-backend/internal/usecase/approval"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type approvalRequestDetailResponse struct {
	ApprovalRequest approvalRequestResponse      `json:"approval_request"`
	Decisions       []approvalDecisionResponse   `json:"decisions"`
	AuditEvents     []approvalAuditEventResponse `json:"audit_events"`
}

type approvalAuditEventResponse struct {
	Type      string    `json:"type"`
	ActorID   string    `json:"actor_id"`
	ActorName string    `json:"actor_name,omitempty"`
	Summary   string    `json:"summary"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// @Summary		Find Approval Request Detail
// @Description	Retrieve an approval request with decision history and audit events
// @Tags			Approval
// @Produce		json
// @Param			approval-request	path		string																		true	"Approval Request ID"
// @Success		200					{object}	httpresponse.SuccessResponse{data=approvalRequestDetailResponse,metadata=nil}	"Approval request found"
// @Failure		400					{object}	httpresponse.ErrorResponse{data=nil}										"Bad request"
// @Failure		404					{object}	httpresponse.ErrorResponse{data=nil}										"Approval request not found"
// @Failure		500					{object}	httpresponse.ErrorResponse{data=nil}										"Internal server error"
// @Router			/approval-requests/{approval-request} [get]
func (h *approvalHandler) FindApprovalRequestDetail(c *gin.Context) {
	output, err := h.approvalUsecase.FindApprovalRequestDetail(c.Request.Context(), approvalUsecase.FindApprovalRequestDetailInput{
		ID: c.Param("approval-request"),
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newApprovalRequestDetailResponse(output))
}

func (h *approvalHandler) newApprovalRequestDetailResponse(output *approvalUsecase.FindApprovalRequestDetailOutput) approvalRequestDetailResponse {
	if output == nil || output.ApprovalRequest == nil {
		return approvalRequestDetailResponse{
			Decisions:   []approvalDecisionResponse{},
			AuditEvents: []approvalAuditEventResponse{},
		}
	}

	decisions := make([]approvalDecisionResponse, 0, len(output.Decisions))
	for _, decision := range output.Decisions {
		response := h.newApprovalDecisionResponse(&decision)
		response.DecidedByName = output.ActorNames[decision.DecidedBy]
		decisions = append(decisions, response)
	}

	auditEvents := make([]approvalAuditEventResponse, 0, len(output.AuditEvents))
	for _, event := range output.AuditEvents {
		auditEvents = append(auditEvents, approvalAuditEventResponse{
			Type:      event.Type,
			ActorID:   event.ActorID.String(),
			ActorName: event.ActorName,
			Summary:   event.Summary,
			Comment:   event.Comment,
			CreatedAt: event.CreatedAt,
		})
	}

	approvalRequest := h.newApprovalRequestResponse(output.ApprovalRequest)
	approvalRequest.ResourceName = output.ResourceName
	approvalRequest.RequestedByName = output.RequestedByName
	approvalRequest.Scope = output.Scope

	return approvalRequestDetailResponse{
		ApprovalRequest: approvalRequest,
		Decisions:       decisions,
		AuditEvents:     auditEvents,
	}
}
