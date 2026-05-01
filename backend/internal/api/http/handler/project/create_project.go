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

type createEnvironmentConfigRequest struct {
	Domain       string `json:"domain" binding:"required"`
	APIDomain    string `json:"api_domain" binding:"required"`
	Region       string `json:"region" binding:"required"`
	IngressClass string `json:"ingress_class,omitempty" binding:"required"`
	PublicAccess bool   `json:"public_access,omitempty" binding:"required"`
}

type createEnvironmentRequest struct {
	Name           string                          `json:"name" binding:"required"` // dev, prod
	Tier           string                          `json:"tier" binding:"required"` // development, production
	Cluster        string                          `json:"cluster" binding:"required"`
	Namespace      string                          `json:"namespace" binding:"required"`
	ArgoCDInstance string                          `json:"argocd_instance" binding:"required"`
	Config         *createEnvironmentConfigRequest `json:"config,omitempty"`
}

type createProjectRequest struct {
	Name         string                     `json:"name" binding:"required"`
	Description  *string                    `json:"description,omitempty"`
	OwnerTeamID  string                     `json:"owner_team_id" binding:"required"`
	Environments []createEnvironmentRequest `json:"environments,omitempty"`
}

type createProjectResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" example:"Project Name"`
	Description string `json:"description" example:"Project Description"`
	OwnerTeamID string `json:"owner_team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedBy   string `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
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
		OwnerTeamID:  input.OwnerTeamID,
		Environments: h.constructEnvironmentInput(input),
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
	return createProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		OwnerTeamID: project.OwnerTeamID.String(),
		CreatedBy:   project.CreatedBy.String(),
	}
}

func (h *projectHandler) constructEnvironmentInput(input createProjectRequest) []projectUsecase.CreateEnvironmentInput {
	environments := make([]projectUsecase.CreateEnvironmentInput, 0, len(input.Environments))
	for _, environment := range input.Environments {
		environments = append(environments, projectUsecase.CreateEnvironmentInput{
			Name:           environment.Name,
			Tier:           environment.Tier,
			Cluster:        environment.Cluster,
			Namespace:      environment.Namespace,
			ArgoCDInstance: environment.ArgoCDInstance,
			Config: func() *projectUsecase.CreateEnvironmentConfigInput {
				if environment.Config == nil {
					return nil
				}
				return &projectUsecase.CreateEnvironmentConfigInput{
					Domain:       environment.Config.Domain,
					APIDomain:    environment.Config.APIDomain,
					Region:       environment.Config.Region,
					IngressClass: environment.Config.IngressClass,
					PublicAccess: environment.Config.PublicAccess,
				}
			}(),
		})
	}

	return environments
}
