package handler

import (
	"devhub-backend/internal/config"
	deploymentUsecase "devhub-backend/internal/usecase/deployment"

	"github.com/gin-gonic/gin"
)

type DeploymentHandler interface {
	CreateDeployment(c *gin.Context)
	FindDeploymentByID(c *gin.Context)
	FindAllDeployments(c *gin.Context)
	UpdateDeployment(c *gin.Context)
	DeleteDeployment(c *gin.Context)
}

type deploymentHandler struct {
	appConfig         config.AppConfig
	deploymentUsecase deploymentUsecase.DeploymentUsecase
}

func NewDeploymentHandler(appConfig config.AppConfig, deploymentUsecase deploymentUsecase.DeploymentUsecase) DeploymentHandler {
	return &deploymentHandler{
		appConfig:         appConfig,
		deploymentUsecase: deploymentUsecase,
	}
}
