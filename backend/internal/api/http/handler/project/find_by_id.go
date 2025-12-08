package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/util/httpresponse"

	projectUsecase "devhub-backend/internal/usecase/project"

	"github.com/gin-gonic/gin"
)

type findOneProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[development, production]"`
	CreatedAt    string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
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

	return findOneProjectResponse{
		ID:           project.ID.String(),
		Description:  project.Description,
		Environments: project.Environments,
		CreatedAt:    project.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    project.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
