package handler

import (
	"devhub-backend/internal/util/httpresponse"

	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"

	"github.com/gin-gonic/gin"
)

// @Summary		Delete Scaffold Request by ID
// @Description	Delete a scaffold request by its ID
// @Tags			Scaffold Request
// @Produce		json
// @Param			scaffold-request	path		string																	true	"Scaffold Request ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=nil,metadata=nil}	"Scaffold Request deleted"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Scaffold Request not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/scaffold-requests/:scaffold-request [delete]
func (h *scaffoldRequestHandler) DeleteScaffoldRequest(c *gin.Context) {
	scaffoldRequestID := c.Param("scaffold-request")
	_, err := h.scaffoldRequestUsecase.DeleteScaffoldRequest(c.Request.Context(), scaffoldRequestUsecase.DeleteScaffoldRequestInput{
		ID: scaffoldRequestID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, nil)
}
