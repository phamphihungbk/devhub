package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createScaffoldRequest struct {
	Template    string            `json:"template" binding:"required"`
	Environment string            `json:"environment" binding:"required,oneof=dev staging prod"`
	Variables   scaffoldVariables `json:"variables" binding:"required,dive"`
}

type scaffoldVariables struct {
	ServiceName   string `json:"service_name" binding:"required"`
	Port          int    `json:"port" binding:"required"`
	Database      string `json:"database" binding:"required,oneof=postgres mysql"`
	EnableLogging bool   `json:"enable_logging"`
}

type createScaffoldRequestResponse struct {
	ID          string                 `json:"id" example:"ad5b0c1f-762a-4ab3-a3e9-50a9057c49f3"`
	Template    string                 `json:"template" example:"go-service"`
	ProjectID   string                 `json:"project_id" example:"1a221b2c-abb7-44c0-8a96-8e92638b2422"`
	Environment string                 `json:"environment" example:"dev"`
	Variables   map[string]interface{} `json:"variables" example:"{\"service_name\":\"payment-service\",\"port\":8080,\"database\":\"postgres\",\"enable_logging\":true}"`
}

// @Summary		Create Scaffold Request
// @Description	Create a new scaffold request
// @Tags			ScaffoldRequest
// @Accept			json
// @Produce		json
// @Param			request	body		createScaffoldRequest													true	"Scaffold request creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createScaffoldRequestResponse,metadata=nil}	"Scaffold request created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/projects/{project}/scaffold_requests [post]
func (h *scaffoldRequestHandler) CreateScaffoldRequest(c *gin.Context) {
	projectID := c.Param("project")
	var input createScaffoldRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdScaffoldRequest, err := h.scaffoldRequestUsecase.CreateScaffoldRequest(c.Request.Context(), scaffoldRequestUsecase.CreateScaffoldRequestInput{
		ProjectID:   projectID,
		Template:    input.Template,
		Environment: input.Environment,
		Variables:   input.Variables,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateScaffoldRequestResponse(createdScaffoldRequest))
}

func (h *scaffoldRequestHandler) newCreateScaffoldRequestResponse(scaffoldRequest *entity.ScaffoldRequest) createScaffoldRequestResponse {
	if scaffoldRequest == nil {
		return createScaffoldRequestResponse{}
	}

	return createScaffoldRequestResponse{
		ID:          scaffoldRequest.ID.String(),
		Name:        scaffoldRequest.Name,
		Description: scaffoldRequest.Description,
		Environments: func() []string {
			envs := make([]string, len(scaffoldRequest.Environments))
			for i, env := range scaffoldRequest.Environments {
				envs[i] = string(env)
			}
			return envs
		}(),
	}
}
