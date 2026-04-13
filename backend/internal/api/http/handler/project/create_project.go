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
	OwnerTeam    string   `json:"owner_team" example:"platform" binding:"required"`
	RepoURL      string   `json:"repo_url" example:"https://git.example.com/acme/project.git" binding:"required"`
	RepoProvider string   `json:"repo_provider" example:"gitea" binding:"required"`
	OwnerContact string   `json:"owner_contact" example:"team@example.com" binding:"required"`
}

type createProjectResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"Project Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[dev, prod, staging]"`
	Status       string   `json:"status" example:"active"`
	OwnerTeam    string   `json:"owner_team" example:"platform"`
	RepoURL      string   `json:"repo_url" example:"https://git.example.com/acme/project.git"`
	RepoProvider string   `json:"repo_provider" example:"gitea"`
	OwnerContact string   `json:"owner_contact" example:"team@example.com"`
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
		OwnerTeam:    input.OwnerTeam,
		RepoURL:      input.RepoURL,
		RepoProvider: input.RepoProvider,
		OwnerContact: input.OwnerContact,
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
		OwnerTeam:    project.OwnerTeam,
		RepoURL:      project.RepoURL,
		RepoProvider: project.RepoProvider,
		OwnerContact: project.OwnerContact,
		CreatedBy:    project.CreatedBy.String(),
	}
}
