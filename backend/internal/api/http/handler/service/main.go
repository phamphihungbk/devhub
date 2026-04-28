package handler

import (
	"devhub-backend/internal/config"
	serviceUsecase "devhub-backend/internal/usecase/service"

	"github.com/gin-gonic/gin"
)

type ServiceHandler interface {
	FindAllServices(c *gin.Context)
	FindServiceDependencies(c *gin.Context)
	CreateServiceDependency(c *gin.Context)
	DeleteServiceDependency(c *gin.Context)
}

type serviceHandler struct {
	appConfig      config.AppConfig
	serviceUsecase serviceUsecase.ServiceUsecase
}

func NewServiceHandler(appConfig config.AppConfig, serviceUsecase serviceUsecase.ServiceUsecase) ServiceHandler {
	return &serviceHandler{
		appConfig:      appConfig,
		serviceUsecase: serviceUsecase,
	}
}
