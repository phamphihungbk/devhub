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
	Description  *string  `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[dev, prod, staging]" binding:"required"`
	Status       string   `json:"status" example:"active" binding:"required"`
	TeamID       string   `json:"team_id" example:"123e4567-e89b-12d3-a456-426614174000" binding:"required"`
	ScmProvider  string   `json:"scm_provider" example:"gitea" binding:"required"`
}

type createProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[dev, prod, staging]"`
	Status       string   `json:"status" example:"active"`
	TeamID       string   `json:"team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ScmProvider  string   `json:"scm_provider" example:"gitea"`
	CreatedBy    string   `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
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
	userID, exists := c.Get("user_id")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
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
		Status:       input.Status,
		TeamID:       input.TeamID,
		ScmProvider:  input.ScmProvider,
		CreatedBy:    userID.(string),
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
	envs := make([]string, 0, len(project.Environments))

	for _, env := range project.Environments {
		envs = append(envs, env.String())
	}

	return createProjectResponse{
		ID:           project.ID.String(),
		Name:         project.Name,
		Description:  project.Description,
		Environments: envs,
		Status:       project.Status.String(),
		TeamID:       project.TeamID.String(),
		ScmProvider:  project.ScmProvider,
		CreatedBy:    project.CreatedBy.String(),
	}
}
