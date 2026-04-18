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

type updateProjectRequest struct {
	Name         *string   `json:"name" example:"Project Name"`
	Description  *string   `json:"description" example:"Project Description"`
	Environments *[]string `json:"environments" example:"[development, production]"`
	Status       *string   `json:"status" example:"active" binding:"required"`
	OwnerTeam    *string   `json:"owner_team" example:"platform" binding:"required"`
	ScmProvider  *string   `json:"scm_provider" example:"gitea" binding:"required"`
	OwnerContact *string   `json:"owner_contact" example:"team@example.com" binding:"required"`
}

type updateProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[development, production]"`
	Status       string   `json:"status" example:"active"`
	OwnerTeam    string   `json:"owner_team" example:"platform"`
	ScmProvider  string   `json:"scm_provider" example:"gitea"`
	OwnerContact string   `json:"owner_contact" example:"team@example.com"`
	CreatedBy    string   `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
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
		ID:           projectID,
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
		Status:       input.Status,
		OwnerTeam:    input.OwnerTeam,
		ScmProvider:  input.ScmProvider,
		OwnerContact: input.OwnerContact,
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
	envs := make([]string, 0, len(project.Environments))

	for _, env := range project.Environments {
		envs = append(envs, env.String())
	}

	return updateProjectResponse{
		ID:           project.ID.String(),
		Name:         project.Name,
		Description:  project.Description,
		Environments: envs,
		Status:       project.Status.String(),
		OwnerTeam:    project.OwnerTeam,
		ScmProvider:  project.ScmProvider,
		OwnerContact: project.OwnerContact,
		CreatedBy:    project.CreatedBy.String(),
	}
}
