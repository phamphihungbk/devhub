package handler

import (
	"devhub-backend/internal/util/httpresponse"

	projectUsecase "devhub-backend/internal/usecase/project"

	"github.com/gin-gonic/gin"
)

type findOneProjectResponse struct {
	ID            string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name          string `json:"name" example:"Project Name"`
	Description   string `json:"description" example:"Project Description"`
	OwnerTeamName string `json:"owner_team_name" example:"Team Dev"`
	CreatedByName string `json:"created_by_name" example:"Hung"`
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

func (h *projectHandler) newFindOneProjectResponse(project *projectUsecase.ProjectDetail) findOneProjectResponse {
	if project == nil {
		return findOneProjectResponse{}
	}

	return findOneProjectResponse{
		ID:            project.ID.String(),
		Name:          project.Name,
		Description:   project.Description,
		OwnerTeamName: project.OwnerTeamName,
		CreatedByName: project.CreatorName,
	}
}
