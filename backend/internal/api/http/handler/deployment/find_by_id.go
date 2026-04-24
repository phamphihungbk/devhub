package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/util/httpresponse"

	deploymentUsecase "devhub-backend/internal/usecase/deployment"

	"github.com/gin-gonic/gin"
)

type findOneDeploymentResponse struct {
	ID           string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceID    string `json:"service_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PluginID     string `json:"plugin_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Environment  string `json:"environment" example:"prod"`
	Version      string `json:"version" example:"1.0.0"`
	Status       string `json:"status" example:"Deployment Status"`
	ExternalRef  string `json:"external_ref" example:"argocd-sync-123"`
	CommitSHA    string `json:"commit_sha" example:"abc123def456"`
	RunnerOutput string `json:"runner_output,omitempty"`
	RunnerError  string `json:"runner_error,omitempty"`
	TriggeredBy  string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary		Find Deployment by ID
// @Description	Retrieve deployment details by its ID
// @Tags			Deployment
// @Produce		json
// @Param			id	path		string																	true	"Deployment ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=findOneDeploymentResponse,metadata=nil}	"Deployment found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Deployment not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/deployment/:deployment [get]
func (h *deploymentHandler) FindDeploymentByID(c *gin.Context) {
	deploymentID := c.Param("deployment")
	deployment, err := h.deploymentUsecase.FindOneDeployment(c.Request.Context(), deploymentUsecase.FindOneDeploymentInput{
		ID: deploymentID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindOneDeploymentResponse(deployment))
}

func (h *deploymentHandler) newFindOneDeploymentResponse(deployment *entity.Deployment) findOneDeploymentResponse {
	if deployment == nil {
		return findOneDeploymentResponse{}
	}

	return findOneDeploymentResponse{
		ID:           deployment.ID.String(),
		ServiceID:    deployment.ServiceID.String(),
		PluginID:     deployment.PluginID.String(),
		Environment:  deployment.Environment.String(),
		Version:      deployment.Version,
		Status:       deployment.Status.String(),
		ExternalRef:  deployment.ExternalRef,
		CommitSHA:    deployment.CommitSHA,
		RunnerOutput: deployment.RunnerOutput,
		RunnerError:  deployment.RunnerError,
		TriggeredBy:  deployment.TriggeredBy.String(),
	}
}
