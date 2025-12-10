package handler

import (
	"devhub-backend/internal/config"
	authUsecase "devhub-backend/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetMe(c *gin.Context)
}

type authHandler struct {
	appConfig   config.AppConfig
	authUsecase authUsecase.AuthUsecase
}

func NewAuthHandler(appConfig config.AppConfig, authUsecase authUsecase.AuthUsecase) AuthHandler {
	return &authHandler{
		appConfig:   appConfig,
		authUsecase: authUsecase,
	}
}
