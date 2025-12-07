package httproute

import (
	userHandler "devhub-backend/internal/api/http/handler/user"
	"devhub-backend/internal/api/http/middleware"
	"devhub-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type Router interface {
	RegisterRoutes(router *gin.Engine)
}

type router struct {
	cfg         config.AppConfig        // Configuration for the application
	Middleware  middleware.Middleware   // Middleware for handling requests
	UserHandler userHandler.UserHandler // Handler for user routes
}

type Dependency struct {
	Middleware  middleware.Middleware
	UserHandler userHandler.UserHandler
}

func NewHTTPRoutes(cfg config.AppConfig, dep Dependency) Router {
	return &router{
		cfg:         cfg,
		Middleware:  dep.Middleware,
		UserHandler: dep.UserHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyUserRoutes(router)
	// r.applyHealthCheckRoutes(router)
	// r.applyConcertRoutes(router)
	// r.applySeatReservationRoutes(router)
}

// applyUserRoutes applies the user routes to the provided router
func (r *router) applyUserRoutes(router *gin.Engine) {
	userRoute := router.Group("/users")
	{
		userRoute.GET("/", r.UserHandler.FindAllUsers)
		userRoute.POST("/", r.UserHandler.CreateUser)
		userRoute.GET("/:id", r.UserHandler.FindUserByID)
		userRoute.DELETE("/:id", r.UserHandler.DeleteUser)
		userRoute.PATCH("/:id", r.UserHandler.UpdateUser)
	}
}

// applyHealthCheckRoutes applies the health check routes to the provided router
// func (r *router) applyHealthCheckRoutes(router *gin.Engine) {
// 	healthRoute := router.Group("/health")
// 	{
// 		healthRoute.GET("/liveness", r.Middleware.BasicAuth(r.cfg.AdminAPIKey, r.cfg.AdminAPISecret), r.HealthCheckHandler.Liveness)
// 		healthRoute.GET("/readiness", r.Middleware.BasicAuth(r.cfg.AdminAPIKey, r.cfg.AdminAPISecret), r.HealthCheckHandler.Readiness)
// 	}
// }

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
