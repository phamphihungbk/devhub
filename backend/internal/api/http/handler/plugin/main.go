package handler

import (
	"devhub-backend/internal/config"
	pluginUsecase "devhub-backend/internal/usecase/plugin"

	"github.com/gin-gonic/gin"
)

type PluginHandler interface {
	FindAllPlugins(c *gin.Context)
	FindPluginByID(c *gin.Context)
	SyncPlugins(c *gin.Context)
	DisablePlugin(c *gin.Context)
}

type pluginHandler struct {
	appConfig     config.AppConfig
	pluginUsecase pluginUsecase.PluginUsecase
}

func NewPluginHandler(appConfig config.AppConfig, pluginUsecase pluginUsecase.PluginUsecase) PluginHandler {
	return &pluginHandler{
		appConfig:     appConfig,
		pluginUsecase: pluginUsecase,
	}
}
