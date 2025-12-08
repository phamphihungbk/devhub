package handler

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type getMeResponse struct {
	ID           string   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name         string   `json:"name" example:"User Name"`
	Description  string   `json:"description" example:"Project Description"`
	Environments []string `json:"environments" example:"[prod,dev,staging]"`
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
	user, exists := c.Get("user")

	if !exists {
		httpresponse.Error(c, errs.NewBadRequestError("unauthorized", nil))
		return
	}

	userEntity, ok := user.(*entity.User)

	if !ok {
		httpresponse.Error(c, errs.NewBadRequestError("invalid user type", nil))
		return
	}

	httpresponse.Success(c, h.newGetMeResponse(userEntity))
}

func (h *authHandler) newGetMeResponse(user *entity.User) getMeResponse {
	if user == nil {
		return getMeResponse{}
	}

	return getMeResponse{
		ID:   user.ID.String(),
		Name: user.Name,
	}
}
