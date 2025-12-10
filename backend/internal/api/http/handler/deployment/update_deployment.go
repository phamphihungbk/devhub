package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	deploymentUsecase "devhub-backend/internal/usecase/deployment"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateDeploymentRequest struct {
	Environment *string `json:"environment" example:"prod"`
	Service     *string `json:"service" example:"Service Name"`
	Version     *string `json:"version" example:"1.0.0"`
	Status      *string `json:"status" example:"Deployment Status"`
}

type updateDeploymentResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ProjectID   string `json:"project_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Environment string `json:"environment" example:"prod"`
	Service     string `json:"service" example:"Service Name"`
	Version     string `json:"version" example:"1.0.0"`
	Status      string `json:"status" example:"Deployment Status"`
	TriggeredBy string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary		Update Deployment
// @Description	Update an existing deployment
// @Tags			Deployment
// @Accept			json
// @Produce		json
// @Param			request	body		updateDeploymentRequest													true	"Deployment update input"
// @Success		200		{object}	httpresponse.SuccessResponse{data=updateDeploymentResponse,metadata=nil}	    "Deployment updated"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/deployments/{deployment} [patch]
func (h *deploymentHandler) UpdateDeployment(c *gin.Context) {
	deploymentID := c.Param("deployment")
	var input updateDeploymentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	updatedDeployment, err := h.deploymentUsecase.UpdateDeployment(c.Request.Context(), deploymentUsecase.UpdateDeploymentInput{
		ID:          deploymentID,
		Environment: input.Environment,
		Service:     input.Service,
		Version:     input.Version,
		Status:      input.Status,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, h.newUpdateDeploymentResponse(updatedDeployment))
}

func (h *deploymentHandler) newUpdateDeploymentResponse(deployment *entity.Deployment) updateDeploymentResponse {
	if deployment == nil {
		return updateDeploymentResponse{}
	}

	return updateDeploymentResponse{
		ID:          deployment.ID.String(),
		ProjectID:   deployment.ProjectID.String(),
		Environment: deployment.Environment.String(),
		Service:     deployment.Service,
		Version:     deployment.Version,
		Status:      deployment.Status.String(),
		TriggeredBy: deployment.TriggeredBy.String(),
	}
}
