package server

import (
	httpmiddleware "devhub-backend/internal/api/http/middleware"
	infraLogger "devhub-backend/internal/infra/logger"

	"github.com/gin-gonic/gin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func (s *Server) setupMiddlewares(appLogger infraLogger.Logger, tracerProvider *sdktrace.TracerProvider) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		httpmiddleware.RequestID(),
		httpmiddleware.Trace(
			httpmiddleware.WithTracerProvider(tracerProvider),
		),
		httpmiddleware.RequestLogger(
			httpmiddleware.WithRequestLogger(appLogger),
		),
	}
}
