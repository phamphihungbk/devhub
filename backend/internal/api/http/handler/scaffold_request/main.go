package handler

import (
	"devhub-backend/internal/config"

	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"

	"github.com/gin-gonic/gin"
)

type ScaffoldRequestHandler interface {
	CreateScaffoldRequest(c *gin.Context)
	SuggestScaffoldRequest(c *gin.Context)
	FindScaffoldRequestByID(c *gin.Context)
	FindAllScaffoldRequests(c *gin.Context)
	DeleteScaffoldRequest(c *gin.Context)
}

type scaffoldRequestHandler struct {
	appConfig              config.AppConfig
	scaffoldRequestUsecase scaffoldRequestUsecase.ScaffoldRequestUsecase
}

func NewScaffoldRequestHandler(appConfig config.AppConfig, scaffoldRequestUsecase scaffoldRequestUsecase.ScaffoldRequestUsecase) ScaffoldRequestHandler {
	return &scaffoldRequestHandler{
		appConfig:              appConfig,
		scaffoldRequestUsecase: scaffoldRequestUsecase,
	}
}
