package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	Auth(tokenSecret string) gin.HandlerFunc
}

type middleware struct{}

func New() Middleware {
	return &middleware{}
}
