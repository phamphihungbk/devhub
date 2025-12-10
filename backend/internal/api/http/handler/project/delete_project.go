package handler

import (
	"devhub-backend/internal/util/httpresponse"

	projectUsecase "devhub-backend/internal/usecase/project"

	"github.com/gin-gonic/gin"
)

// @Summary		Delete Project by ID
// @Description	Delete a project by its ID
// @Tags			Project
// @Produce		json
// @Param			project	path		string																	true	"Project ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=nil,metadata=nil}	"Project deleted"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"User not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/{project} [delete]
func (h *projectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("project")

	_, err := h.projectUsecase.DeleteProject(c.Request.Context(), projectUsecase.DeleteProjectInput{
		ID: projectID,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, nil)
}
