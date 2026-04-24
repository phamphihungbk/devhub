package middleware

import (
	"devhub-backend/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Auth(tokenSecret string) gin.HandlerFunc
	RequirePermissions(permissions ...entity.Permission) gin.HandlerFunc
	Authorize(tokenSecret string, permissions ...entity.Permission) gin.HandlersChain
}

type middleware struct{}

var _ Middleware = (*middleware)(nil)

func New() Middleware {
	return &middleware{}
}

func (m *middleware) Authorize(tokenSecret string, permissions ...entity.Permission) gin.HandlersChain {
	handlers := gin.HandlersChain{
		m.Auth(tokenSecret),
	}

	if len(permissions) > 0 {
		handlers = append(handlers, m.RequirePermissions(permissions...))
	}

	return handlers
}
