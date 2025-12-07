package handler

import (
	"devhub-backend/internal/config"
	userUsecase "devhub-backend/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	FindUserByID(c *gin.Context)
	FindAllUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	appConfig   config.AppConfig
	userUsecase userUsecase.UserUsecase
}

func NewUserHandler(appConfig config.AppConfig, userUsecase userUsecase.UserUsecase) UserHandler {
	return &userHandler{
		appConfig:   appConfig,
		userUsecase: userUsecase,
	}
}
