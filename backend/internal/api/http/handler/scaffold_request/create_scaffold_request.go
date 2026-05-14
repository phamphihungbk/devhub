package handler

import (
	"devhub-backend/internal/api/http/approvaltarget"
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createScaffoldRequest struct {
	PluginID      string                 `json:"plugin_id" binding:"required"`
	EnvironmentID string                 `json:"environment_id" binding:"required"`
	Variables     map[string]interface{} `json:"variables" binding:"required"`
}

type createScaffoldRequestResponse struct {
	ID          string                          `json:"id" example:"ad5b0c1f-762a-4ab3-a3e9-50a9057c49f3"`
	RequestedBy string                          `json:"requested_by" example:"8bb6438e-b4a7-4945-9969-f446f7c26ca5"`
	Status      string                          `json:"status" example:"pending"`
	ProjectID   string                          `json:"project_id" example:"1a221b2c-abb7-44c0-8a96-8e92638b2422"`
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

	approvalResource, _ := approvaltarget.StringsFromContext(c)

	usecaseInput := scaffoldRequestUsecase.CreateScaffoldRequestInput{
		PluginID:         input.PluginID,
		EnvironmentID:    input.EnvironmentID,
		ProjectID:        projectID,
		RequestedBy:      userID.(string),
		Variables:        input.Variables,
		ApprovalResource: approvalResource,
	}

	createdScaffoldRequest, err := h.scaffoldRequestUsecase.CreateScaffoldRequest(c.Request.Context(), usecaseInput)

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
		RequestedBy: scaffoldRequest.RequestedBy.String(),
		Status:      scaffoldRequest.Status.String(),
		ProjectID:   scaffoldRequest.ProjectID.String(),
		Variables:   scaffoldRequest.Variables,
	}
}
