package handler

import (
	"devhub-backend/internal/config"
	projectUsecase "devhub-backend/internal/usecase/project"

	"github.com/gin-gonic/gin"
)

type ProjectHandler interface {
	CreateProject(c *gin.Context)
	FindProjectByID(c *gin.Context)
	FindAllProjects(c *gin.Context)
	UpdateProject(c *gin.Context)
	DeleteProject(c *gin.Context)
}

type projectHandler struct {
	appConfig      config.AppConfig
	projectUsecase projectUsecase.ProjectUsecase
}

func NewProjectHandler(appConfig config.AppConfig, projectUsecase projectUsecase.ProjectUsecase) ProjectHandler {
	return &projectHandler{
		appConfig:      appConfig,
		projectUsecase: projectUsecase,
	}
}
