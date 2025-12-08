package handler

import (
	"devhub-backend/internal/config"
	pluginUsecase "devhub-backend/internal/usecase/plugin"

	"github.com/gin-gonic/gin"
)

type PluginHandler interface {
	CreatePlugin(c *gin.Context)
	FindPluginByID(c *gin.Context)
	FindAllPlugins(c *gin.Context)
	UpdatePlugin(c *gin.Context)
	DeletePlugin(c *gin.Context)
}

type pluginHandler struct {
	appConfig     config.AppConfig
	pluginUsecase pluginUsecase.PluginUsecase
}

func NewPluginHandler(appConfig config.AppConfig, pluginUsecase pluginUsecase.PluginUsecase) PluginHandler {
	return &projectHandler{
		appConfig:      appConfig,
		projectUsecase: projectUsecase,
	}
}
