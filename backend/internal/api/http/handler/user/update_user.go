package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	userUsecase "devhub-backend/internal/usecase/user"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateUserRequest struct {
	Name *string `json:"name" example:"User Name" binding:"required"`
	Role *string `json:"role" example:"user" binding:"required" validate:"oneof=admin user"`
}

type updateUserResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"User Name"`
	Email string `json:"email" example:"user@example.com"`
}

// @Summary		Update User
// @Description	Update an existing user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	path		string
// @Param			request	body		updateUserRequest													true	"User update input"
// @Success		200		{object}	httpresponse.SuccessResponse{data=updateUserResponse,metadata=nil}	    "User updated"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/users/{user} [patch]
func (h *userHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("user")
	var input updateUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	updatedUser, err := h.userUsecase.UpdateUser(c.Request.Context(), userUsecase.UpdateUserInput{
		ID:   userID,
		Name: input.Name,
		Role: input.Role,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusOK, h.newUpdateUserResponse(updatedUser))
}

func (h *userHandler) newUpdateUserResponse(user *entity.User) updateUserResponse {
	if user == nil {
		return updateUserResponse{}
	}

	return updateUserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}
}
