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
	Environment string `json:"environment" example:"prod" binding:"required"`
	Service     string `json:"service" example:"Service Name" binding:"required"`
	Version     string `json:"version" example:"v1.0.0" binding:"required"`
}

type createDeploymentResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ProjectID   string `json:"project_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Environment string `json:"environment" example:"prod"`
	Service     string `json:"service" example:"Service Name"`
	Version     string `json:"version" example:"1.0.0"`
	Status      string `json:"status" example:"Deployment Status"`
	ExternalRef string `json:"external_ref" example:"argocd-sync-123"`
	CommitSHA   string `json:"commit_sha" example:"abc123def456"`
	TriggeredBy string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000"`
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
// @Router			/projects/:project/deployments [post]
func (h *deploymentHandler) CreateDeployment(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	projectID := c.Param("project")
	var input createDeploymentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdDeployment, err := h.deploymentUsecase.CreateDeployment(c.Request.Context(), deploymentUsecase.CreateDeploymentInput{
		ProjectID:   projectID,
		Environment: input.Environment,
		Service:     input.Service,
		Version:     input.Version,
		TriggeredBy: userID.(string),
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
		ID:          deployment.ID.String(),
		ProjectID:   deployment.ProjectID.String(),
		Environment: deployment.Environment.String(),
		Service:     deployment.Service,
		Version:     deployment.Version,
		Status:      deployment.Status.String(),
		ExternalRef: deployment.ExternalRef,
		CommitSHA:   deployment.CommitSHA,
		TriggeredBy: deployment.TriggeredBy.String(),
	}
}
