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

type createUserRequest struct {
	Name     *string `json:"name" example:"User Name"`
	Email    string  `json:"email" example:"user@example.com" binding:"required"`
	Password string  `json:"password" example:"password123" binding:"required"`
	Role     string  `json:"role" example:"user" binding:"required"`
}

type createUserResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"User Name"`
	Email string `json:"email" example:"user@example.com"`
	Role  string `json:"role" example:"admin"`
}

// @Summary		Create User
// @Description	Create a new user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			request	body		createUserRequest													true	"User creation input"
// @Success		201		{object}	httpresponse.SuccessResponse{data=createUserResponse,metadata=nil}	"User created"
// @Failure		400		{object}	httpresponse.ErrorResponse{data=nil}									"Bad request"
// @Failure		500		{object}	httpresponse.ErrorResponse{data=nil}									"Internal server error"
// @Router			/users [post]
func (h *userHandler) CreateUser(c *gin.Context) {
	var input createUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdUser, err := h.userUsecase.CreateUser(c.Request.Context(), userUsecase.CreateUserInput{
		Name:     input.Name,
		Email:    input.Email,
		Role:     input.Role,
		Password: input.Password,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateUserResponse(createdUser))
}

func (h *userHandler) newCreateUserResponse(user *entity.User) createUserResponse {
	if user == nil {
		return createUserResponse{}
	}

	return createUserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.String(),
	}
}
