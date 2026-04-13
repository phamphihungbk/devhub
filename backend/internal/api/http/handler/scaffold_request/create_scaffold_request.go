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

type scaffoldVariablesRequest struct {
	ServiceName   string `json:"service_name" binding:"required"`
	Port          int    `json:"port" binding:"required"`
	Database      string `json:"database" binding:"required"`
	EnableLogging bool   `json:"enable_logging" binding:"required"`
}

type createScaffoldRequest struct {
	PluginID    string                   `json:"plugin_id" binding:"required"`
	Template    string                   `json:"template" binding:"required"`
	Environment string                   `json:"environment" binding:"required"`
	Variables   scaffoldVariablesRequest `json:"variables" binding:"required"`
}

type createScaffoldRequestResponse struct {
	ID          string                          `json:"id" example:"ad5b0c1f-762a-4ab3-a3e9-50a9057c49f3"`
	PluginID    string                          `json:"plugin_id" example:"72bd5b8f-54b3-442a-b54f-685643f6d46e"`
	RequestedBy string                          `json:"requested_by" example:"8bb6438e-b4a7-4945-9969-f446f7c26ca5"`
	Template    string                          `json:"template" example:"go-http@v2"`
	Status      string                          `json:"status" example:"pending"`
	ProjectID   string                          `json:"project_id" example:"1a221b2c-abb7-44c0-8a96-8e92638b2422"`
	Environment string                          `json:"environment" example:"dev"`
	Variables   entity.ScaffoldRequestVariables `json:"variables" example:"{\"service_name\":\"payment-service\",\"port\":8080,\"database\":\"postgres\",\"enable_logging\":true}"`
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
	userID, exists := c.Get("user_id")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	projectID := c.Param("project")
	var input createScaffoldRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdScaffoldRequest, err := h.scaffoldRequestUsecase.CreateScaffoldRequest(c.Request.Context(), scaffoldRequestUsecase.CreateScaffoldRequestInput{
		PluginID:    input.PluginID,
		ProjectID:   projectID,
		RequestedBy: userID.(string),
		Template:    input.Template,
		Environment: input.Environment,
		Variables:   scaffoldRequestUsecase.ScaffoldRequestVariables(input.Variables),
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
		PluginID:    scaffoldRequest.PluginID.String(),
		RequestedBy: scaffoldRequest.RequestedBy.String(),
		Template:    scaffoldRequest.Template,
		Status:      scaffoldRequest.Status.String(),
		ProjectID:   scaffoldRequest.ProjectID.String(),
		Environment: scaffoldRequest.Environment.String(),
		Variables:   scaffoldRequest.Variables,
	}
}
