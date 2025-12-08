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

type createDeploymentRequest struct {
	Name         string   `json:"name" example:"Project Name" binding:"required"`
	Description  string   `json:"description" example:"Project Description" binding:"required"`
	Environments []string `json:"environments" example:"[prod,dev,staging]" binding:"required,dive,required"`
}

type createDeploymentResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[prod,dev,staging]"`
}

// @Summary		Create Deployment
// @Description	Create a new deployment
// @Tags			Deployment
// @Accept			json
// @Produce		json
// @Param			request	body		createDeploymentRequest													true	"Deployment creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createDeploymentResponse,metadata=nil}	"Deployment created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/:project/deployment [post]
func (h *deploymentHandler) CreateDeployment(c *gin.Context) {
	projectID := c.Param("project")
	var input createDeploymentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdDeployment, err := h.deploymentUsecase.CreateDeployment(c.Request.Context(), deploymentUsecase.CreateDeploymentInput{
		ProjectID:    projectID,
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateDeploymentResponse(createdDeployment))
}

func (h *deploymentHandler) newCreateDeploymentResponse(deployment *entity.Deployment) createDeploymentResponse {
	if deployment == nil {
		return createDeploymentResponse{}
	}

	return createDeploymentResponse{
		ID:           deployment.ID.String(),
		Name:         deployment.Name,
		Description:  deployment.Description,
		Environments: deployment.Environments,
	}
}
