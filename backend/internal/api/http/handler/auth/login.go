package handler

import (
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	authUsecase "devhub-backend/internal/usecase/auth"

	"net/http"

	"github.com/gin-gonic/gin"
)

type loginUserRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required"`
	Password string `json:"password" example:"password123" binding:"required"`
}

type loginUserResponse struct {
	AccessToken  string `json:"access_token" example:"jwt-token-string"`
	RefreshToken string `json:"refresh_token" example:"refresh-token-string"`
}

// @Summary		Login
// @Description	Authenticate user and return a JWT token
// @Tags			Auth
// @Produce		json
// @Param			id	path		string																	true	"Project ID"
// @Success		200	{object}	httpresponse.SuccessResponse{data=loginUserResponse,metadata=nil}	"Login successful"
// @Failure		400	{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		404	{object}	httpresponse.ErrorResponse{data=nil}									"User not found"
// @Failure		500	{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var input loginUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	token, err := h.authUsecase.IssueToken(c.Request.Context(), authUsecase.IssueTokenInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, loginUserResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}
