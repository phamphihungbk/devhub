package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/util/httpresponse"

	projectUsecase "devhub-backend/internal/usecase/project"

	"github.com/gin-gonic/gin"
)

type findOneProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[development, production]"`
	Status       string   `json:"status" example:"active"`
	OwnerTeam    string   `json:"owner_team" example:"platform"`
	ScmProvider  string   `json:"scm_provider" example:"gitea"`
	OwnerContact string   `json:"owner_contact" example:"team@example.com"`
	CreatedBy    string   `json:"created_by" example:"Hung Pham"`
}

// @Summary		Find Project by ID
// @Description	Retrieve project details by its ID
// @Tags			Project
// @Produce		json
// @Param			project	path		string																	true	"Project ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=findOneProjectResponse,metadata=nil}	"Project found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Concert not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/{project} [get]
func (h *projectHandler) FindProjectByID(c *gin.Context) {
	projectID := c.Param("project")
	project, err := h.projectUsecase.FindOneProject(c.Request.Context(), projectUsecase.FindOneProjectInput{
		ID: projectID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindOneProjectResponse(project))
}

func (h *projectHandler) newFindOneProjectResponse(project *entity.Project) findOneProjectResponse {
	if project == nil {
		return findOneProjectResponse{}
	}
	envs := make([]string, 0, len(project.Environments))

	for _, env := range project.Environments {
		envs = append(envs, env.String())
	}

	return findOneProjectResponse{
		ID:           project.ID.String(),
		Name:         project.Name,
		Description:  project.Description,
		Environments: envs,
		Status:       project.Status.String(),
		OwnerTeam:    project.OwnerTeam,
		ScmProvider:  project.ScmProvider,
		OwnerContact: project.OwnerContact,
		CreatedBy:    project.CreatedByName,
	}
}
