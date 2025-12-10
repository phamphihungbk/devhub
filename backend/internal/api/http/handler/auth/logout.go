package handler

import (
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/httpresponse"

	authUsecase "devhub-backend/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

// @Summary		Find Project by ID
// @Description	Retrieve project details by its ID
// @Tags			Project
// @Produce		json
// @Param			id	path		string																	true	"Project ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=findOneProjectResponse,metadata=nil}	"Project found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Project not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	_, err := h.authUsecase.RevokeToken(c.Request.Context(), authUsecase.RevokeTokenInput{
		UserID: userID.(string),
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, nil)
}
