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

type createProjectRequest struct {
	Name         string   `json:"name" example:"Project Name" binding:"required"`
	Description  string   `json:"description" example:"Project Description" binding:"required"`
	Environments []string `json:"environments" example:"[development, production]" binding:"required"`
}

type createProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[development, production]"`
	CreatedAt    string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// @Summary		Create Project
// @Description	Create a new project
// @Tags			Project
// @Accept			json
// @Produce		json
// @Param			request	body		createProjectRequest													true	"Project creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createProjectResponse,metadata=nil}	"Project created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects [post]
func (h *projectHandler) CreateProject(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	userEntity, ok := user.(*entity.User)

	if !ok {
		httpresponse.Error(c, errs.NewBadRequestError("invalid user type", nil))
		return
	}

	var input createProjectRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdProject, err := h.projectUsecase.CreateProject(c.Request.Context(), projectUsecase.CreateProjectInput{
		Name:         input.Name,
		Description:  input.Description,
		Environments: input.Environments,
		CreatedBy:    userEntity.ID,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateProjectResponse(createdProject))
}

func (h *projectHandler) newCreateProjectResponse(project *entity.Project) createProjectResponse {
	if project == nil {
		return createProjectResponse{}
	}

	return createProjectResponse{
		ID:           project.ID.String(),
		Name:         project.Name,
		Description:  project.Description,
		Environments: project.Environments,
		CreatedAt:    project.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    project.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
