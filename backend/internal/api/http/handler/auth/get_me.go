package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	authUsecase "devhub-backend/internal/usecase/auth"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type getMeResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"User Name"`
	Email string `json:"email" example:"user@example.com"`
	Role  string `json:"role" example:"admin"`
}

// @Summary		Return Authenticated User Info
// @Description	Retrieve authenticated user details
// @Tags			User
// @Produce		json
// @Success		200	{object}	httpresponse.SuccessResponse{data=getMeResponse,metadata=nil}	"User found"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"User not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/auth/me [get]
func (h *authHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	user, err := h.authUsecase.FindOneUser(c.Request.Context(), authUsecase.FindOneUserInput{
		ID: userID.(string),
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newGetMeResponse(user))
}

func (h *authHandler) newGetMeResponse(user *entity.User) getMeResponse {
	if user == nil {
		return getMeResponse{}
	}

	return getMeResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.String(),
	}
}
