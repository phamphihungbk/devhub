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
	Name         *string                     `json:"name" example:"Project Name"`
	Description  *string                     `json:"description" example:"Project Description"`
	OwnerTeamID  *string                     `json:"owner_team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Environments *[]createEnvironmentRequest `json:"environments,omitempty"`
}

type updateProjectResponse struct {
	ID           string                          `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string                          `json:"name" example:"Project Name"`
	Description  string                          `json:"description" example:"Project Description"`
	OwnerTeamID  string                          `json:"owner_team_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Environments []findOneProjectEnvironmentItem `json:"environments"`
	CreatedBy    string                          `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
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
		Name:         input.Name,
		Description:  input.Description,
		OwnerTeamID:  input.OwnerTeamID,
		Environments: h.constructUpdateEnvironmentInput(input.Environments),
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
		Name:         project.Name,
		Description:  project.Description,
		OwnerTeamID:  project.OwnerTeamID.String(),
		Environments: h.constructEnvironmentResponse(project.Environments),
		CreatedBy:    project.CreatedBy.String(),
	}
}

func (h *projectHandler) constructUpdateEnvironmentInput(input *[]createEnvironmentRequest) *[]projectUsecase.CreateEnvironmentInput {
	if input == nil {
		return nil
	}

	environments := make([]projectUsecase.CreateEnvironmentInput, 0, len(*input))
	for _, environment := range *input {
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

	return &environments
}
