package middleware

import (
	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Auth(tokenSecret string) gin.HandlerFunc
	RequirePermissions(permissions ...entity.PermissionName) gin.HandlerFunc
	Authorize(tokenSecret string, permissions ...entity.PermissionName) gin.HandlersChain
}

type middleware struct {
	userRepository repository.UserRepository
}

func New(userRepository repository.UserRepository) Middleware {
	return &middleware{userRepository: userRepository}
}

func (m *middleware) Authorize(tokenSecret string, permissions ...entity.PermissionName) gin.HandlersChain {
	handlers := gin.HandlersChain{
		m.Auth(tokenSecret),
	}

	if len(permissions) > 0 {
		handlers = append(handlers, m.RequirePermissions(permissions...))
	}

	return handlers
}
