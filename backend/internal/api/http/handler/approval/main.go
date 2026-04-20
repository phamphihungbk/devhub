package handler

import (
	"devhub-backend/internal/config"
	approvalUsecase "devhub-backend/internal/usecase/approval"

	"github.com/gin-gonic/gin"
)

type ApprovalHandler interface {
	CreateApprovalPolicy(c *gin.Context)
	CreateApprovalDecision(c *gin.Context)
}

type approvalHandler struct {
	appConfig       config.AppConfig
	approvalUsecase approvalUsecase.ApprovalUsecase
}

func NewApprovalHandler(appConfig config.AppConfig, approvalUsecase approvalUsecase.ApprovalUsecase) ApprovalHandler {
	return &approvalHandler{
		appConfig:       appConfig,
		approvalUsecase: approvalUsecase,
	}
}
