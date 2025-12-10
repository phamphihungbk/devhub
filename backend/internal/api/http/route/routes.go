package httproute

import (
	authHandler "devhub-backend/internal/api/http/handler/auth"
	deploymentHandler "devhub-backend/internal/api/http/handler/deployment"
	pluginHandler "devhub-backend/internal/api/http/handler/plugin"
	projectHandler "devhub-backend/internal/api/http/handler/project"
	scaffoldRequestHandler "devhub-backend/internal/api/http/handler/scaffold_request"
	userHandler "devhub-backend/internal/api/http/handler/user"
	"devhub-backend/internal/api/http/middleware"
	"devhub-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type Router interface {
	RegisterRoutes(router *gin.Engine)
}

type router struct {
	appCfg                 config.AppConfig                              // Configuration for the application
	tokenCfg               config.TokenConfig                            // Configuration for the application
	Middleware             middleware.Middleware                         // Middleware for handling requests
	UserHandler            userHandler.UserHandler                       // Handler for user routes
	AuthHandler            authHandler.AuthHandler                       // Handler for auth routes
	PluginHandler          pluginHandler.PluginHandler                   // Handler for plugin request routes
	DeploymentHandler      deploymentHandler.DeploymentHandler           // Handler for deployment routes
	ProjectHandler         projectHandler.ProjectHandler                 // Handler for project routes
	ScaffoldRequestHandler scaffoldRequestHandler.ScaffoldRequestHandler // Handler for scaffold request routes
}

type Dependency struct {
	Middleware             middleware.Middleware
	UserHandler            userHandler.UserHandler
	AuthHandler            authHandler.AuthHandler
	PluginHandler          pluginHandler.PluginHandler
	DeploymentHandler      deploymentHandler.DeploymentHandler
	ProjectHandler         projectHandler.ProjectHandler
	ScaffoldRequestHandler scaffoldRequestHandler.ScaffoldRequestHandler
}

func NewHTTPRoutes(appCfg config.AppConfig, tokenCfg config.TokenConfig, dep Dependency) Router {
	return &router{
		appCfg:                 appCfg,
		tokenCfg:               tokenCfg,
		Middleware:             dep.Middleware,
		UserHandler:            dep.UserHandler,
		AuthHandler:            dep.AuthHandler,
		ScaffoldRequestHandler: dep.ScaffoldRequestHandler,
		PluginHandler:          dep.PluginHandler,
		DeploymentHandler:      dep.DeploymentHandler,
		ProjectHandler:         dep.ProjectHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyAuthRoutes(router)
	r.applyUserRoutes(router)
	r.applyProjectRoutes(router)
	r.applyScaffoldRequestRoutes(router)
	r.applyDeploymentRoutes(router)
	r.applyPluginRoutes(router)
}

// applyAuthRoutes applies the auth routes to the provided router
func (r *router) applyAuthRoutes(router *gin.Engine) {
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/login", r.AuthHandler.Login)
		authRoute.POST("/logout", r.AuthHandler.Logout)
		authRoute.GET("/me", r.Middleware.Auth(r.tokenCfg.Secret), r.AuthHandler.GetMe)
	}
}

// applyUserRoutes applies the user routes to the provided router
func (r *router) applyUserRoutes(router *gin.Engine) {
	userRoute := router.Group("/users")
	{
		userRoute.GET("/", r.UserHandler.FindAllUsers)
		userRoute.POST("/", r.UserHandler.CreateUser)
		userRoute.GET("/:user", r.UserHandler.FindUserByID)
		userRoute.DELETE("/:user", r.UserHandler.DeleteUser)
		userRoute.PATCH("/:user", r.UserHandler.UpdateUser)
	}
}

// applyProjectRoutes applies the project routes to the provided router
func (r *router) applyProjectRoutes(router *gin.Engine) {
	projectRoute := router.Group("/projects")
	{
		projectRoute.GET("/", r.ProjectHandler.FindAllProjects)
		projectRoute.POST("/", r.Middleware.Auth(r.tokenCfg.Secret), r.ProjectHandler.CreateProject)
		projectRoute.GET("/:project", r.ProjectHandler.FindProjectByID)
		projectRoute.DELETE("/:project", r.ProjectHandler.DeleteProject)
		projectRoute.PATCH("/:project", r.ProjectHandler.UpdateProject)
	}
}

// applyScaffoldRequestRoutes applies the scaffold request routes to the provided router
func (r *router) applyScaffoldRequestRoutes(router *gin.Engine) {
	scaffoldRequestRoute := router.Group("/scaffold-requests")
	projectRoute := router.Group("/projects/:project")
	{
		projectRoute.GET("/scaffold-requests", r.ScaffoldRequestHandler.FindAllScaffoldRequests)
		projectRoute.POST("/scaffold-requests", r.ScaffoldRequestHandler.CreateScaffoldRequest)

		scaffoldRequestRoute.GET("/:scaffold-request", r.ScaffoldRequestHandler.FindScaffoldRequestByID)
		scaffoldRequestRoute.DELETE("/:scaffold-request", r.ScaffoldRequestHandler.DeleteScaffoldRequest)
	}
}

// applyDeploymentRoutes applies the deployment routes to the provided router
func (r *router) applyDeploymentRoutes(router *gin.Engine) {
	projectRoute := router.Group("/projects/:project")
	deploymentRoute := router.Group("/deployments")
	{
		projectRoute.GET("/deployments", r.DeploymentHandler.FindAllDeployments)
		projectRoute.POST("/deployments", r.DeploymentHandler.CreateDeployment)

		deploymentRoute.GET("/:deployment", r.DeploymentHandler.FindDeploymentByID)
		deploymentRoute.DELETE("/:deployment", r.DeploymentHandler.DeleteDeployment)
		deploymentRoute.PATCH("/:deployment", r.DeploymentHandler.UpdateDeployment)
	}
}

// applyPluginRoutes applies the plugin routes to the provided router
func (r *router) applyPluginRoutes(router *gin.Engine) {
	pluginRoute := router.Group("/plugins")
	{
		pluginRoute.GET("/", r.PluginHandler.FindAllPlugins)
		pluginRoute.POST("/", r.PluginHandler.CreatePlugin)
		pluginRoute.GET("/:id", r.PluginHandler.FindPluginByID)
		pluginRoute.DELETE("/:id", r.PluginHandler.DeletePlugin)
		pluginRoute.PATCH("/:id", r.PluginHandler.UpdatePlugin)
	}
}
