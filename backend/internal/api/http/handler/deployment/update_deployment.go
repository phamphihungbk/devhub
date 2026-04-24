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
	Environment  *string `json:"environment" example:"prod"`
	Version      *string `json:"version" example:"1.0.0"`
	Status       *string `json:"status" example:"Deployment Status"`
	ExternalRef  *string `json:"external_ref" example:"argocd-sync-123"`
	CommitSHA    *string `json:"commit_sha" example:"abc123def456"`
	RunnerOutput *string `json:"runner_output"`
	RunnerError  *string `json:"runner_error"`
	FinishedAt   *string `json:"finished_at" example:"2026-04-13T12:34:56Z"`
}

type updateDeploymentResponse struct {
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
	FinishedAt   string `json:"finished_at,omitempty" example:"2026-04-13T12:34:56Z"`
	TriggeredBy  string `json:"triggered_by" example:"123e4567-e89b-12d3-a456-426614174000"`
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
		ID:           deploymentID,
		Environment:  input.Environment,
		Version:      input.Version,
		Status:       input.Status,
		ExternalRef:  input.ExternalRef,
		CommitSHA:    input.CommitSHA,
		RunnerOutput: input.RunnerOutput,
		RunnerError:  input.RunnerError,
		FinishedAt:   input.FinishedAt,
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
		FinishedAt: func() string {
			if deployment.FinishedAt == nil {
				return ""
			}
			return deployment.FinishedAt.Format("2006-01-02T15:04:05Z07:00")
		}(),
		TriggeredBy: deployment.TriggeredBy.String(),
	}
}
