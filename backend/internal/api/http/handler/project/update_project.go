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
	Name         string    `json:"name" example:"Project Name" binding:"required"`
	Description  *string   `json:"description" example:"Project Description"`
	Environments *[]string `json:"environments" example:"[development, production]"`
}

type updateProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[development, production]"`
	CreatedAt    string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
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
		Description:  input.Description,
		Environments: input.Environments,
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
		ID:           project.ID.String(),
		Description:  project.Description,
		Environments: project.Environments,
		CreatedAt:    project.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    project.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
