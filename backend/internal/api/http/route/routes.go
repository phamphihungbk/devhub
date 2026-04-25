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
	teamHandler "devhub-backend/internal/api/http/handler/team"
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
	TeamHandler            teamHandler.TeamHandler                       // Handler for team routes
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
	TeamHandler            teamHandler.TeamHandler
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
		TeamHandler:            dep.TeamHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyAuthRoutes(router)
	r.applyApprovalRoutes(router)
	r.applyUserRoutes(router)
	r.applyTeamRoutes(router)
	r.applyProjectRoutes(router)
	r.applyScaffoldRequestRoutes(router)
	r.applyDeploymentRoutes(router)
	r.applyPluginRoutes(router)
}

func (r *router) applyTeamRoutes(router *gin.Engine) {
	teamRoute := router.Group("/teams")
	teamRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		teamRoute.GET("/",
			r.Middleware.RequirePermissions(entity.PermissionUserRead),
			r.TeamHandler.FindAllTeams,
		)
		teamRoute.POST("/",
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.TeamHandler.CreateTeam,
		)
		teamRoute.PATCH("/:team",
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.TeamHandler.UpdateTeam,
		)
	}
}

func (r *router) applyApprovalRoutes(router *gin.Engine) {
	approvalPolicyRoute := router.Group("/approval-policies")
	approvalRequestRoute := router.Group("/approval-requests")
	approvalPolicyRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	approvalRequestRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		approvalRequestRoute.GET("/",
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ApprovalHandler.FindAllApprovalRequests,
		)
		approvalRequestRoute.GET("/:approval-request",
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ApprovalHandler.FindApprovalRequestDetail,
		)

		approvalPolicyRoute.POST("/",
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ApprovalHandler.CreateApprovalPolicy,
		)

		approvalRequestRoute.POST("/:approval-request/decisions",
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ApprovalHandler.CreateApprovalDecision,
		)
	}
}

// applyAuthRoutes applies the auth routes to the provided router
func (r *router) applyAuthRoutes(router *gin.Engine) {
	authRoute := router.Group("/auth")
	authProtectedRoute := authRoute.Group("/")
	authProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		authRoute.POST("/login", r.AuthHandler.Login)
		authProtectedRoute.POST("/logout", r.AuthHandler.Logout)
		authProtectedRoute.GET("/me", r.AuthHandler.GetMe)
	}
}

// applyUserRoutes applies the user routes to the provided router
func (r *router) applyUserRoutes(router *gin.Engine) {
	userRoute := router.Group("/users")
	userRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		userRoute.GET("/",
			r.Middleware.RequirePermissions(entity.PermissionUserRead),
			r.UserHandler.FindAllUsers,
		)
		userRoute.POST("/",
			r.Middleware.RequirePermissions(entity.PermissionUserWrite),
			r.UserHandler.CreateUser,
		)
		userRoute.GET("/:user",
			r.Middleware.RequirePermissions(entity.PermissionUserRead),
			r.UserHandler.FindUserByID,
		)
		userRoute.DELETE("/:user",
			r.Middleware.RequirePermissions(entity.PermissionUserWrite),
			r.UserHandler.DeleteUser,
		)
		userRoute.PATCH("/:user",
			r.Middleware.RequirePermissions(entity.PermissionUserWrite),
			r.UserHandler.UpdateUser,
		)
	}
}

// applyProjectRoutes applies the project routes to the provided router
func (r *router) applyProjectRoutes(router *gin.Engine) {
	projectRoute := router.Group("/projects")
	projectProtectedRoute := projectRoute.Group("/")
	projectProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		projectRoute.GET("/", r.ProjectHandler.FindAllProjects)
		projectProtectedRoute.POST("/",
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ProjectHandler.CreateProject,
		)
		projectRoute.GET("/:project", r.ProjectHandler.FindProjectByID)
		projectRoute.GET("/:project/services", r.ServiceHandler.FindAllServices)
		projectProtectedRoute.DELETE("/:project",
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ProjectHandler.DeleteProject,
		)
		projectProtectedRoute.PATCH("/:project",
			r.Middleware.RequirePermissions(entity.PermissionProjectWrite),
			r.ProjectHandler.UpdateProject,
		)
	}

	serviceRoute := router.Group("/services")
	serviceProtectedRoute := serviceRoute.Group("/")
	serviceProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		serviceRoute.GET("/:service/deployments", r.DeploymentHandler.FindAllDeployments)
		serviceProtectedRoute.POST("/:service/deployments",
			r.Middleware.RequirePermissions(entity.PermissionDeploymentWrite),
			r.DeploymentHandler.CreateDeployment,
		)
		serviceRoute.GET("/:service/releases", r.ReleaseHandler.FindAllReleases)
		serviceProtectedRoute.POST("/:service/releases",
			r.Middleware.RequirePermissions(entity.PermissionReleaseWrite),
			r.ReleaseHandler.CreateRelease,
		)
		serviceProtectedRoute.POST("/:service/scaffold-suggestions",
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ServiceHandler.SuggestScaffold,
		)
	}
}

// applyScaffoldRequestRoutes applies the scaffold request routes to the provided router
func (r *router) applyScaffoldRequestRoutes(router *gin.Engine) {
	scaffoldRequestRoute := router.Group("/scaffold-requests")
	projectRoute := router.Group("/projects/:project")
	projectProtectedRoute := projectRoute.Group("/")
	projectProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	scaffoldRequestProtectedRoute := scaffoldRequestRoute.Group("/")
	scaffoldRequestProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		projectRoute.GET("/scaffold-requests", r.ScaffoldRequestHandler.FindAllScaffoldRequests)
		projectProtectedRoute.POST("/scaffold-requests",
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ScaffoldRequestHandler.CreateScaffoldRequest,
		)

		scaffoldRequestRoute.GET("/:scaffold-request", r.ScaffoldRequestHandler.FindScaffoldRequestByID)
		scaffoldRequestProtectedRoute.DELETE("/:scaffold-request",
			r.Middleware.RequirePermissions(entity.PermissionScaffoldRequestWrite),
			r.ScaffoldRequestHandler.DeleteScaffoldRequest,
		)
	}
}

// applyDeploymentRoutes applies the deployment routes to the provided router
func (r *router) applyDeploymentRoutes(router *gin.Engine) {
	deploymentRoute := router.Group("/deployments")
	deploymentProtectedRoute := deploymentRoute.Group("/")
	deploymentProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		deploymentRoute.GET("/:deployment", r.DeploymentHandler.FindDeploymentByID)
		deploymentProtectedRoute.DELETE("/:deployment",
			r.Middleware.RequirePermissions(entity.PermissionDeploymentWrite),
			r.DeploymentHandler.DeleteDeployment,
		)
		deploymentProtectedRoute.PATCH("/:deployment",
			r.Middleware.RequirePermissions(entity.PermissionDeploymentWrite),
			r.DeploymentHandler.UpdateDeployment,
		)
	}
}

// applyPluginRoutes applies the plugin routes to the provided router
func (r *router) applyPluginRoutes(router *gin.Engine) {
	pluginRoute := router.Group("/plugins")
	pluginProtectedRoute := pluginRoute.Group("/")
	pluginProtectedRoute.Use(r.Middleware.Auth(r.tokenCfg.Secret))
	{
		pluginRoute.GET("/", r.PluginHandler.FindAllPlugins)
		pluginProtectedRoute.POST("/",
			r.Middleware.RequirePermissions(entity.PermissionPluginWrite),
			r.PluginHandler.CreatePlugin,
		)
		pluginRoute.GET("/:plugin", r.PluginHandler.FindPluginByID)
		pluginProtectedRoute.DELETE("/:plugin",
			r.Middleware.RequirePermissions(entity.PermissionPluginWrite),
			r.PluginHandler.DeletePlugin,
		)
		pluginProtectedRoute.PATCH("/:plugin",
			r.Middleware.RequirePermissions(entity.PermissionPluginWrite),
			r.PluginHandler.UpdatePlugin,
		)
	}
}
