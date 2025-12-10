package handler

import (
	"devhub-backend/internal/util/httpresponse"

	deploymentUsecase "devhub-backend/internal/usecase/deployment"

	"github.com/gin-gonic/gin"
)

// @Summary		Delete Deployment by ID
// @Description	Delete a deployment by its ID
// @Tags			Deployment
// @Produce		json
// @Param			id	path		string																	true	"Deployment ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=nil,metadata=nil}	"Deployment deleted"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Deployment not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/deployments/{deployment} [delete]
func (h *deploymentHandler) DeleteDeployment(c *gin.Context) {
	deploymentID := c.Param("deployment")
	_, err := h.deploymentUsecase.DeleteDeployment(c.Request.Context(), deploymentUsecase.DeleteDeploymentInput{
		ID: deploymentID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, nil)
}
