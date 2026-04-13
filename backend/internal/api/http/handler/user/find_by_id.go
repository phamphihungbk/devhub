package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/util/httpresponse"

	userUsecase "devhub-backend/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type findOneUserResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"User Name"`
	Email string `json:"email" example:"user@example.com"`
	Role  string `json:"role" example:"platform_admin"`
}

// @Summary		Find User by ID
// @Description	Retrieve user details by its ID
// @Tags			User
// @Produce		json
// @Param			user	path		string																	true	"User ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=findOneUserResponse,metadata=nil}	"User found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"Concert not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/users/{user} [get]
func (h *userHandler) FindUserByID(c *gin.Context) {
	userID := c.Param("user")
	user, err := h.userUsecase.FindOneUser(c.Request.Context(), userUsecase.FindOneUserInput{
		ID: userID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindOneUserResponse(user))
}

func (h *userHandler) newFindOneUserResponse(user *entity.User) findOneUserResponse {
	if user == nil {
		return findOneUserResponse{}
	}

	return findOneUserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.String(),
	}
}
