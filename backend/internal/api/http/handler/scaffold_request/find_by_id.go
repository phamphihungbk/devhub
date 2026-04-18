package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/util/httpresponse"

	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"

	"github.com/gin-gonic/gin"
)

type findOneScaffoldRequestResponse struct {
	ID           string                          `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PluginID     string                          `json:"plugin_id" example:"72bd5b8f-54b3-442a-b54f-685643f6d46e"`
	RequestedBy  string                          `json:"requested_by" example:"8bb6438e-b4a7-4945-9969-f446f7c26ca5"`
	ProjectID    string                          `json:"project_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Status       string                          `json:"status" example:"pending"`
	Environments string                          `json:"environments" example:"prod"`
	Variables    entity.ScaffoldRequestVariables `json:"variables" example:"key=value"`
}

// @Summary		Find Scaffold Request by ID
// @Description	Retrieve scaffold request details by its ID
// @Tags			Scaffold Request
// @Produce		json
// @Param			scaffold-request	path		string																	true	"Scaffold Request ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=findOneScaffoldRequestResponse,metadata=nil}	"Scaffold Request found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Scaffold Request not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/scaffold-requests/{scaffold-request} [get]
func (h *scaffoldRequestHandler) FindScaffoldRequestByID(c *gin.Context) {
	scaffoldRequestID := c.Param("scaffold-request")
	scaffoldRequest, err := h.scaffoldRequestUsecase.FindOneScaffoldRequest(c.Request.Context(), scaffoldRequestUsecase.FindOneScaffoldRequestInput{
		ID: scaffoldRequestID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindOneScaffoldRequestResponse(scaffoldRequest))
}

func (h *scaffoldRequestHandler) newFindOneScaffoldRequestResponse(scaffoldRequest *entity.ScaffoldRequest) findOneScaffoldRequestResponse {
	if scaffoldRequest == nil {
		return findOneScaffoldRequestResponse{}
	}

	return findOneScaffoldRequestResponse{
		ID:           scaffoldRequest.ID.String(),
		PluginID:     scaffoldRequest.PluginID.String(),
		RequestedBy:  scaffoldRequest.RequestedBy.String(),
		ProjectID:    scaffoldRequest.ProjectID.String(),
		Status:       scaffoldRequest.Status.String(),
		Environments: scaffoldRequest.Environment.String(),
		Variables:    scaffoldRequest.Variables,
	}
}
