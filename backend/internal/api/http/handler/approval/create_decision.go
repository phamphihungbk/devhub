package handler

import (
	"net/http"
	"strings"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	approvalUsecase "devhub-backend/internal/usecase/approval"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type createApprovalDecisionRequest struct {
	Decision string `json:"decision" binding:"required"`
	Comment  string `json:"comment"`
}

type createApprovalDecisionResponse struct {
	ApprovalRequest approvalRequestResponse  `json:"approval_request"`
	Decision        approvalDecisionResponse `json:"decision"`
}

type approvalRequestResponse struct {
	ID                string     `json:"id"`
	Resource          string     `json:"resource"`
	ResourceName      string     `json:"resource_name,omitempty"`
	Action            string     `json:"action"`
	ResourceID        string     `json:"resource_id"`
	RequestedBy       string     `json:"requested_by"`
	RequestedByName   string     `json:"requested_by_name,omitempty"`
	Scope             string     `json:"scope,omitempty"`
	ProjectID         string     `json:"project_id,omitempty"`
	ServiceID         string     `json:"service_id,omitempty"`
	Environment       string     `json:"environment,omitempty"`
	Status            string     `json:"status"`
	RequiredApprovals int        `json:"required_approvals"`
	ApprovedCount     int        `json:"approved_count"`
	RejectedCount     int        `json:"rejected_count"`
	ResolvedAt        *time.Time `json:"resolved_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type approvalDecisionResponse struct {
	DecidedBy     string    `json:"decided_by"`
	DecidedByName string    `json:"decided_by_name,omitempty"`
	Decision      string    `json:"decision"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
}

// @Summary		Create Approval Decision
// @Description	Create an approval decision for an approval request
// @Tags			Approval
// @Accept			json
// @Produce		json
// @Param			approval-request	path		string																		true	"Approval Request ID"
// @Param			request				body		createApprovalDecisionRequest												true	"Approval decision creation input"
// @Success		201					{object}	httpresponse.SuccessResponse{data=createApprovalDecisionResponse,metadata=nil}	"Approval decision created"
// @Failure		400					{object}	httpresponse.ErrorResponse{data=nil}										"Bad request"
// @Failure		409					{object}	httpresponse.ErrorResponse{data=nil}										"Conflict"
// @Failure		500					{object}	httpresponse.ErrorResponse{data=nil}										"Internal server error"
// @Router			/approval-requests/{approval-request}/decisions [post]
func (h *approvalHandler) CreateApprovalDecision(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	var input createApprovalDecisionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	output, err := h.approvalUsecase.CreateApprovalDecision(c.Request.Context(), approvalUsecase.CreateApprovalDecisionInput{
		ApprovalRequestID: c.Param("approval-request"),
		DecidedBy:         userID.(string),
		Decision:          input.Decision,
		Comment:           input.Comment,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateApprovalDecisionResponse(output))
}

func (h *approvalHandler) newCreateApprovalDecisionResponse(output *approvalUsecase.CreateApprovalDecisionOutput) createApprovalDecisionResponse {
	if output == nil || output.ApprovalRequest == nil || output.Decision == nil {
		return createApprovalDecisionResponse{}
	}

	return createApprovalDecisionResponse{
		ApprovalRequest: h.newApprovalRequestResponse(output.ApprovalRequest),
		Decision:        h.newApprovalDecisionResponse(output.Decision),
	}
}

func (h *approvalHandler) newApprovalRequestResponse(request *entity.ApprovalRequest) approvalRequestResponse {
	if request == nil {
		return approvalRequestResponse{}
	}

	response := approvalRequestResponse{
		ID:                request.ID.String(),
		Resource:          request.Resource,
		Action:            request.Action,
		ResourceID:        request.ResourceID.String(),
		RequestedBy:       request.RequestedBy.String(),
		Status:            request.Status.String(),
		RequiredApprovals: request.RequiredApprovals,
		ApprovedCount:     request.ApprovedCount,
		RejectedCount:     request.RejectedCount,
		ResolvedAt:        request.ResolvedAt,
		CreatedAt:         request.CreatedAt,
		UpdatedAt:         request.UpdatedAt,
	}
	if request.ProjectID != nil {
		response.ProjectID = request.ProjectID.String()
	}
	if request.ServiceID != nil {
		response.ServiceID = request.ServiceID.String()
	}
	if request.Environment != nil {
		response.Environment = *request.Environment
	}
	response.Scope = newApprovalRequestScope(response)

	return response
}

func newApprovalRequestScope(request approvalRequestResponse) string {
	parts := make([]string, 0, 3)
	if request.ProjectID != "" {
		parts = append(parts, request.ProjectID)
	}
	if request.ServiceID != "" {
		parts = append(parts, request.ServiceID)
	}
	if request.Environment != "" {
		parts = append(parts, request.Environment)
	}
	if len(parts) == 0 {
		return "Global"
	}
	return strings.Join(parts, " / ")
}

func (h *approvalHandler) newApprovalDecisionResponse(decision *entity.ApprovalDecision) approvalDecisionResponse {
	if decision == nil {
		return approvalDecisionResponse{}
	}

	return approvalDecisionResponse{
		DecidedBy: decision.DecidedBy.String(),
		Decision:  decision.Decision.String(),
		Comment:   decision.Comment,
		CreatedAt: decision.CreatedAt,
	}
}
