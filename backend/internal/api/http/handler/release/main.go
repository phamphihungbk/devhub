package handler

import (
	"devhub-backend/internal/config"
	releaseUsecase "devhub-backend/internal/usecase/release"

	"github.com/gin-gonic/gin"
)

type ReleaseHandler interface {
	CreateRelease(c *gin.Context)
	FindAllReleases(c *gin.Context)
}

type releaseHandler struct {
	appConfig      config.AppConfig
	releaseUsecase releaseUsecase.ReleaseUsecase
}

func NewReleaseHandler(appConfig config.AppConfig, releaseUsecase releaseUsecase.ReleaseUsecase) ReleaseHandler {
	return &releaseHandler{
		appConfig:      appConfig,
		releaseUsecase: releaseUsecase,
	}
}
