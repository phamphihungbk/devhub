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
	Name         *string   `json:"name" example:"Project Name" binding:"required"`
	Description  *string   `json:"description" example:"Project Description" binding:"required"`
	Environments *[]string `json:"environments" example:"[prod,dev,staging]" binding:"required"`
}

type updateProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[prod,dev,staging]"`
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
// @Router			/projects/{id} [patch]
func (h *projectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")
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

	environments := make([]string, len(project.Environments))
	for i, env := range project.Environments {
		environments[i] = string(env)
	}

	return updateProjectResponse{
		ID:           project.ID.String(),
		Name:         project.Name,
		Description:  project.Description,
		Environments: environments,
	}
}
