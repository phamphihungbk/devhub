package handler

import (
	"devhub-backend/internal/config"
	teamUsecase "devhub-backend/internal/usecase/team"

	"github.com/gin-gonic/gin"
)

type TeamHandler interface {
	CreateTeam(c *gin.Context)
	UpdateTeam(c *gin.Context)
	FindAllTeams(c *gin.Context)
}

type teamHandler struct {
	appConfig   config.AppConfig
	teamUsecase teamUsecase.TeamUsecase
}

func NewTeamHandler(appConfig config.AppConfig, teamUsecase teamUsecase.TeamUsecase) TeamHandler {
	return &teamHandler{
		appConfig:   appConfig,
		teamUsecase: teamUsecase,
	}
}
