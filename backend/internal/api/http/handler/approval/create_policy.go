package handler

import (
	"net/http"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	approvalUsecase "devhub-backend/internal/usecase/approval"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

type createApprovalPolicyRequest struct {
	Resource          string  `json:"resource" binding:"required"`
	ServiceID         *string `json:"service_id"`
	EnvironmentID     string  `json:"environment_id" binding:"required"`
	RequiredApprovals int     `json:"required_approvals" binding:"required"`
	Enabled           bool    `json:"enabled" binding:"required"`
}

type createApprovalPolicyResponse struct {
	ID                string `json:"id"`
	Resource          string `json:"resource"`
	ServiceID         string `json:"service_id,omitempty"`
	EnvironmentID     string `json:"environment_id"`
	RequiredApprovals int    `json:"required_approvals"`
	Enabled           bool   `json:"enabled"`
}

// @Summary		Create Approval Policy
// @Description	Create a new approval policy
// @Tags			Approval
// @Accept			json
// @Produce		json
// @Param			request	body		createApprovalPolicyRequest													true	"Approval policy creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createApprovalPolicyResponse,metadata=nil}	"Approval policy created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}										"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}										"Internal server error"
// @Router			/approval-policies [post]
func (h *approvalHandler) CreateApprovalPolicy(c *gin.Context) {
	var input createApprovalPolicyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	policy, err := h.approvalUsecase.CreateApprovalPolicy(c.Request.Context(), approvalUsecase.CreateApprovalPolicyInput{
		Resource:          input.Resource,
		ServiceID:         input.ServiceID,
		EnvironmentID:     input.EnvironmentID,
		RequiredApprovals: input.RequiredApprovals,
		Enabled:           input.Enabled,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateApprovalPolicyResponse(policy))
}

func (h *approvalHandler) newCreateApprovalPolicyResponse(policy *entity.ApprovalPolicy) createApprovalPolicyResponse {
	if policy == nil {
		return createApprovalPolicyResponse{}
	}

	response := createApprovalPolicyResponse{
		ID:                policy.ID.String(),
		Resource:          policy.Resource,
		EnvironmentID:     policy.EnvironmentID.String(),
		RequiredApprovals: policy.RequiredApprovals,
		Enabled:           policy.Enabled,
	}

	return response
}
