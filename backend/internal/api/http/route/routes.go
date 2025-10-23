package httproute

import (
	"github.com/gin-gonic/gin"
	// "github.com/phamphihungbk/devhub-backend/internal/http/middleware"
)

type Router interface {
	RegisterRoutes(router *gin.Engine)
}

type router struct {
	cfg                config.AppConfig                 // Configuration for the application
	Middleware         middleware.Middleware            // Middleware for handling requests
	HealthCheckHandler healthHandler.HealthCheckHandler // Handler for health check routes
	ConcertHandler     concertHandler.ConcertHandler    // Handler for concert routes
	SeatHandler        seatHandler.SeatHandler          // Handler for seat routes
}

type Dependency struct {
	// Middleware         middleware.Middleware
	// HealthCheckHandler healthHandler.HealthCheckHandler
	// ConcertHandler     concertHandler.ConcertHandler
	// SeatHandler        seatHandler.SeatHandler
}

func NewHTTPRoutes(cfg config.AppConfig, dep Dependency) Router {
	return &router{
		cfg:                cfg,
		Middleware:         dep.Middleware,
		HealthCheckHandler: dep.HealthCheckHandler,
		ConcertHandler:     dep.ConcertHandler,
		SeatHandler:        dep.SeatHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyHealthCheckRoutes(router)
	r.applyConcertRoutes(router)
	r.applySeatReservationRoutes(router)
}

// applyHealthCheckRoutes applies the health check routes to the provided router
func (r *router) applyHealthCheckRoutes(router *gin.Engine) {
	healthRoute := router.Group("/health")
	{
		healthRoute.GET("/liveness", r.Middleware.BasicAuth(r.cfg.AdminAPIKey, r.cfg.AdminAPISecret), r.HealthCheckHandler.Liveness)
		healthRoute.GET("/readiness", r.Middleware.BasicAuth(r.cfg.AdminAPIKey, r.cfg.AdminAPISecret), r.HealthCheckHandler.Readiness)
	}
}

// func RegisterRoutes(r *gin.Engine) {
// 	// r.Use(middleware.Logger())
// 	// r.Use(middleware.Auth())

// 	v1 := r.Group("/api/v1")

// 	// // Service Scaffolding
// 	// services := v1.Group("/services")
// 	// services.POST("", handler.CreateServiceHandler)
// 	// services.GET("/templates", handler.ListServiceTemplatesHandler)
// 	// services.GET("", handler.ListServicesHandler)
// 	// services.GET("/:id", handler.GetServiceHandler)

// 	// // Deployment Management
// 	// deployments := v1.Group("/deployments")
// 	// deployments.POST("", handler.CreateDeploymentHandler)
// 	// deployments.GET("", handler.ListDeploymentsHandler)
// 	// deployments.GET("/:id", handler.GetDeploymentHandler)
// 	// deployments.POST("/:id/rollback", handler.RollbackDeploymentHandler)

// 	// // CI/CD Integration
// 	// cicd := v1.Group("/cicd")
// 	// cicd.GET("/pipelines", handler.ListPipelinesHandler)
// 	// cicd.POST("/pipelines", handler.TriggerPipelineHandler)

// 	// // Metrics & Monitoring
// 	// metrics := v1.Group("/metrics")
// 	// metrics.GET("", handler.ListMetricsHandler)
// 	// metrics.GET("/:serviceId", handler.GetServiceMetricsHandler)

// 	// // Authentication & Access Control
// 	// auth := v1.Group("/auth")
// 	// auth.POST("/login", handler.LoginHandler)
// 	// auth.POST("/logout", handler.LogoutHandler)
// 	// auth.GET("/roles", handler.ListRolesHandler)
// 	// auth.POST("/roles", handler.CreateRoleHandler)

// 	// // Git Integration
// 	// git := v1.Group("/git")
// 	// git.GET("/repos", handler.ListReposHandler)
// 	// git.POST("/repos", handler.ConnectRepoHandler)

// 	// // Background Jobs
// 	// jobs := v1.Group("/jobs")
// 	// jobs.POST("", handler.CreateJobHandler)
// 	// jobs.GET("", handler.ListJobsHandler)
// 	// jobs.GET("/:id", handler.GetJobHandler)

// 	// // Plugin System
// 	// plugins := v1.Group("/plugins")
// 	// plugins.GET("", handler.ListPluginsHandler)
// 	// plugins.POST("", handler.InstallPluginHandler)

// 	// // API Testing
// 	// v1.POST("/api-test", handler.RunApiTestHandler)

// 	// Miscellaneous
// 	v1.GET("/health", handler.HealthHandler)
// }
