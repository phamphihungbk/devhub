package handler

import (
	"devhub-backend/internal/util/httpresponse"

	userUsecase "devhub-backend/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

// @Summary		Delete User by ID
// @Description	Delete a user by its ID
// @Tags			User
// @Produce		json
// @Param			id	path		string																	true	"User ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=nil,metadata=nil}	"User deleted"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"User not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/users/{id} [delete]
func (h *userHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	_, err := h.userUsecase.DeleteUser(c.Request.Context(), userUsecase.DeleteUserInput{
		ID: userID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, nil)
}
