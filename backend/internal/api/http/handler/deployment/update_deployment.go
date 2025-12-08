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
	Name         *string   `json:"name" example:"Deployment Name" binding:"required"`
	Description  *string   `json:"description" example:"Deployment Description" binding:"required"`
	Environments *[]string `json:"environments" example:"[prod,dev,staging]" binding:"required"`
}

type updateDeploymentResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Deployment Name"`
	Description  string   `json:"description" example:"Deployment Description"`
	Environments []string `json:"environments" example:"[prod,dev,staging]"`
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
// @Router			/deployments/{id} [patch]
func (h *deploymentHandler) UpdateDeployment(c *gin.Context) {
	deploymentID := c.Param("id")
	var input updateDeploymentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	updatedDeployment, err := h.deploymentUsecase.UpdateDeployment(c.Request.Context(), deploymentUsecase.UpdateDeploymentInput{
		ID:           deploymentID,
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
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

	environments := make([]string, len(deployment.Environments))
	for i, env := range deployment.Environments {
		environments[i] = string(env)
	}

	return updateDeploymentResponse{
		ID:           deployment.ID.String(),
		Name:         deployment.Name,
		Description:  deployment.Description,
		Environments: environments,
	}
}
