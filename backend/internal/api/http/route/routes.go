package httproute

import (
	approvalHandler "devhub-backend/internal/api/http/handler/approval"
	authHandler "devhub-backend/internal/api/http/handler/auth"
	deploymentHandler "devhub-backend/internal/api/http/handler/deployment"
	pluginHandler "devhub-backend/internal/api/http/handler/plugin"
	projectHandler "devhub-backend/internal/api/http/handler/project"
	releaseHandler "devhub-backend/internal/api/http/handler/release"
	scaffoldRequestHandler "devhub-backend/internal/api/http/handler/scaffold_request"
	serviceHandler "devhub-backend/internal/api/http/handler/service"
	userHandler "devhub-backend/internal/api/http/handler/user"
	"devhub-backend/internal/api/http/middleware"
	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type Router interface {
	RegisterRoutes(router *gin.Engine)
}

type router struct {
	appCfg                 config.AppConfig                              // Configuration for the application
	tokenCfg               config.TokenConfig                            // Configuration for the application
	Middleware             middleware.Middleware                         // Middleware for handling requests
	ApprovalHandler        approvalHandler.ApprovalHandler               // Handler for approval routes
	UserHandler            userHandler.UserHandler                       // Handler for user routes
	AuthHandler            authHandler.AuthHandler                       // Handler for auth routes
	PluginHandler          pluginHandler.PluginHandler                   // Handler for plugin request routes
	DeploymentHandler      deploymentHandler.DeploymentHandler           // Handler for deployment routes
	ProjectHandler         projectHandler.ProjectHandler                 // Handler for project routes
	ReleaseHandler         releaseHandler.ReleaseHandler                 // Handler for release routes
	ScaffoldRequestHandler scaffoldRequestHandler.ScaffoldRequestHandler // Handler for scaffold request routes
	ServiceHandler         serviceHandler.ServiceHandler                 // Handler for service routes
}

type Dependency struct {
	Middleware             middleware.Middleware
	ApprovalHandler        approvalHandler.ApprovalHandler
	UserHandler            userHandler.UserHandler
	AuthHandler            authHandler.AuthHandler
	PluginHandler          pluginHandler.PluginHandler
	DeploymentHandler      deploymentHandler.DeploymentHandler
	ProjectHandler         projectHandler.ProjectHandler
	ReleaseHandler         releaseHandler.ReleaseHandler
	ScaffoldRequestHandler scaffoldRequestHandler.ScaffoldRequestHandler
	ServiceHandler         serviceHandler.ServiceHandler
}

func NewHTTPRoutes(appCfg config.AppConfig, tokenCfg config.TokenConfig, dep Dependency) Router {
	return &router{
		appCfg:                 appCfg,
		tokenCfg:               tokenCfg,
		Middleware:             dep.Middleware,
		ApprovalHandler:        dep.ApprovalHandler,
		UserHandler:            dep.UserHandler,
		AuthHandler:            dep.AuthHandler,
		ScaffoldRequestHandler: dep.ScaffoldRequestHandler,
		PluginHandler:          dep.PluginHandler,
		DeploymentHandler:      dep.DeploymentHandler,
		ProjectHandler:         dep.ProjectHandler,
		ReleaseHandler:         dep.ReleaseHandler,
		ServiceHandler:         dep.ServiceHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyAuthRoutes(router)
	r.applyApprovalRoutes(router)
	r.applyUserRoutes(router)
	r.applyProjectRoutes(router)
	r.applyScaffoldRequestRoutes(router)
	r.applyDeploymentRoutes(router)
	r.applyPluginRoutes(router)
}

func (r *router) applyApprovalRoutes(router *gin.Engine) {
	approvalPolicyRoute := router.Group("/approval-policies")
	approvalRequestRoute := router.Group("/approval-requests")
	{
		approvalPolicyRoute.POST("/",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ApprovalHandler.CreateApprovalPolicy,
		)

		approvalRequestRoute.POST("/:approval-request/decisions",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ApprovalHandler.CreateApprovalDecision,
		)
	}
}

// applyAuthRoutes applies the auth routes to the provided router
func (r *router) applyAuthRoutes(router *gin.Engine) {
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/login", r.AuthHandler.Login)
		authRoute.POST("/logout", r.Middleware.Auth(r.tokenCfg.Secret), r.AuthHandler.Logout)
		authRoute.GET("/me", r.Middleware.Auth(r.tokenCfg.Secret), r.AuthHandler.GetMe)
	}
}

// applyUserRoutes applies the user routes to the provided router
func (r *router) applyUserRoutes(router *gin.Engine) {
	userRoute := router.Group("/users")
	{
		userRoute.GET("/",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionUserRead),
			r.UserHandler.FindAllUsers,
		)
		userRoute.POST("/",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionUserWrite),
			r.UserHandler.CreateUser,
		)
		userRoute.GET("/:user",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionUserRead),
			r.UserHandler.FindUserByID,
		)
		userRoute.DELETE("/:user",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionUserWrite),
			r.UserHandler.DeleteUser,
		)
		userRoute.PATCH("/:user",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionUserWrite),
			r.UserHandler.UpdateUser,
		)
	}
}

// applyProjectRoutes applies the project routes to the provided router
func (r *router) applyProjectRoutes(router *gin.Engine) {
	projectRoute := router.Group("/projects")
	{
		projectRoute.GET("/", r.ProjectHandler.FindAllProjects)
		projectRoute.POST("/",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ProjectHandler.CreateProject,
		)
		projectRoute.GET("/:project", r.ProjectHandler.FindProjectByID)
		projectRoute.GET("/:project/services", r.ServiceHandler.FindAllServices)
		projectRoute.DELETE("/:project",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ProjectHandler.DeleteProject,
		)
		projectRoute.PATCH("/:project",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ProjectHandler.UpdateProject,
		)
	}

	serviceRoute := router.Group("/services")
	{
		serviceRoute.GET("/:service/deployments", r.DeploymentHandler.FindAllDeployments)
		serviceRoute.POST("/:service/deployments",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionDeploymentWrite),
			r.DeploymentHandler.CreateDeployment,
		)
		serviceRoute.GET("/:service/releases", r.ReleaseHandler.FindAllReleases)
		serviceRoute.POST("/:service/releases",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionReleaseWrite),
			r.ReleaseHandler.CreateRelease,
		)
	}
}

// applyScaffoldRequestRoutes applies the scaffold request routes to the provided router
func (r *router) applyScaffoldRequestRoutes(router *gin.Engine) {
	scaffoldRequestRoute := router.Group("/scaffold-requests")
	projectRoute := router.Group("/projects/:project")
	{
		projectRoute.GET("/scaffold-requests", r.ScaffoldRequestHandler.FindAllScaffoldRequests)
		projectRoute.POST("/scaffold-requests",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ScaffoldRequestHandler.CreateScaffoldRequest,
		)

		scaffoldRequestRoute.GET("/:scaffold-request", r.ScaffoldRequestHandler.FindScaffoldRequestByID)
		scaffoldRequestRoute.DELETE("/:scaffold-request",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ScaffoldRequestHandler.DeleteScaffoldRequest,
		)
	}
}

// applyDeploymentRoutes applies the deployment routes to the provided router
func (r *router) applyDeploymentRoutes(router *gin.Engine) {
	deploymentRoute := router.Group("/deployments")
	{
		deploymentRoute.GET("/:deployment", r.DeploymentHandler.FindDeploymentByID)
		deploymentRoute.DELETE("/:deployment",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionDeploymentWrite),
			r.DeploymentHandler.DeleteDeployment,
		)
		deploymentRoute.PATCH("/:deployment",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionDeploymentWrite),
			r.DeploymentHandler.UpdateDeployment,
		)
	}
}

// applyPluginRoutes applies the plugin routes to the provided router
func (r *router) applyPluginRoutes(router *gin.Engine) {
	pluginRoute := router.Group("/plugins")
	{
		pluginRoute.GET("/", r.PluginHandler.FindAllPlugins)
		pluginRoute.POST("/",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionPluginWrite),
			r.PluginHandler.CreatePlugin,
		)
		pluginRoute.GET("/:plugin", r.PluginHandler.FindPluginByID)
		pluginRoute.DELETE("/:plugin",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionPluginWrite),
			r.PluginHandler.DeletePlugin,
		)
		pluginRoute.PATCH("/:plugin",
			r.Middleware.Auth(r.tokenCfg.Secret),
			r.Middleware.RequirePermissions(entity.PermissionPluginWrite),
			r.PluginHandler.UpdatePlugin,
		)
	}
}
