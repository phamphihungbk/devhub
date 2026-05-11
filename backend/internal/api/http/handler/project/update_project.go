package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	projectUsecase "devhub-backend/internal/usecase/project"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateEnvironmentConfigRequest struct {
	Domain       string `json:"domain" binding:"required"`
	APIDomain    string `json:"api_domain" binding:"required"`
	Region       string `json:"region" binding:"required"`
	IngressClass string `json:"ingress_class" binding:"required"`
	PublicAccess bool   `json:"public_access"`
}

type updateEnvironmentRequest struct {
	Name           string                         `json:"name" binding:"required"` // dev, prod
	Tier           string                         `json:"tier" binding:"required"` // development, production
	Cluster        string                         `json:"cluster" binding:"required"`
	Namespace      string                         `json:"namespace" binding:"required"`
	ArgoCDInstance string                         `json:"argocd_instance" binding:"required"`
	Config         updateEnvironmentConfigRequest `json:"config" binding:"required"`
}

type updateProjectRequest struct {
	Name         *string                     `json:"name" example:"Project Name"`
	Description  *string                     `json:"description" example:"Project Description"`
	OwnerTeamID  *string                     `json:"owner_team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Environments *[]updateEnvironmentRequest `json:"environments,omitempty"`
}

type updateProjectResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" example:"Project Name"`
	Description string `json:"description" example:"Project Description"`
	OwnerTeamID string `json:"owner_team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedBy   string `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary		Update Project
// @Description	Update an existing project
// @Tags			Project
// @Accept			json
// @Produce		json
// @Param			request	body		updateProjectRequest													true	"Project update input"
// @Success		200		{object}	httpresponse.SuccessResponse{data=updateProjectResponse,metadata=nil}	    "Project updated"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/{project} [patch]
func (h *projectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("project")
	var input updateProjectRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	updatedProject, err := h.projectUsecase.UpdateProject(c.Request.Context(), projectUsecase.UpdateProjectInput{
		ID:          projectID,
		Name:        input.Name,
		Description: input.Description,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, h.newUpdateProjectResponse(updatedProject))
}

func (h *projectHandler) newUpdateProjectResponse(project *entity.Project) updateProjectResponse {
	if project == nil {
		return updateProjectResponse{}
	}
	return updateProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		OwnerTeamID: project.OwnerTeamID.String(),
		CreatedBy:   project.CreatedBy.String(),
	}
}
