package middleware

import (
	"devhub-backend/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Auth(tokenSecret string) gin.HandlerFunc
	RequirePermissions(permissions ...entity.Permission) gin.HandlerFunc
}

type middleware struct{}

func New() Middleware {
	return &middleware{}
}
