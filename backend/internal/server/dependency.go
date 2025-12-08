package server

import (
	"context"

	"github.com/jmoiron/sqlx"

	authHandler "devhub-backend/internal/api/http/handler/auth"
	deploymentHandler "devhub-backend/internal/api/http/handler/deployment"
	pluginHandler "devhub-backend/internal/api/http/handler/plugin"
	projectHandler "devhub-backend/internal/api/http/handler/project"
	scaffoldRequestHandler "devhub-backend/internal/api/http/handler/scaffold_request"
	userHandler "devhub-backend/internal/api/http/handler/user"
	"devhub-backend/internal/api/http/middleware"
	httproute "devhub-backend/internal/api/http/route"
	dbDeploymentRepo "devhub-backend/internal/infra/db/repository/deployment"
	dbPluginRepo "devhub-backend/internal/infra/db/repository/plugin"
	dbProjectRepo "devhub-backend/internal/infra/db/repository/project"
	dbScaffoldRequestRepo "devhub-backend/internal/infra/db/repository/scaffold_request"
	dbUserRepo "devhub-backend/internal/infra/db/repository/user"
	"devhub-backend/internal/infra/logger"
	deploymentUsecase "devhub-backend/internal/usecase/deployment"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	projectUsecase "devhub-backend/internal/usecase/project"
	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"
	userUsecase "devhub-backend/internal/usecase/user"
)

//nolint:unparam
func (s *Server) setupRouteDependencies(ctx context.Context, appLogger logger.Logger, dbConn *sqlx.DB) (httproute.Dependency, error) {
	// Transactor factory
	// transactorFactory := infraDB.NewSqlxTransactorFactory(dbConn)

	// DB Repositories
	dbUserRepo := dbUserRepo.NewUserRepository(dbConn)
	dbProjectRepo := dbProjectRepo.NewProjectRepository(dbConn)
	dbDeploymentRepo := dbDeploymentRepo.NewDeploymentRepository(dbConn)
	dbPluginRepo := dbPluginRepo.NewPluginRepository(dbConn)
	dbScaffoldRequestRepo := dbScaffoldRequestRepo.NewScaffoldRequestRepository(dbConn)

	// Query retrier
	// queryBackoff, _ := retry.NewExponentialBackoffStrategy(500*time.Millisecond, 2.0, 5*time.Second)
	// queryRetrier, _ := retry.NewRetrier(retry.Config{
	// MaxAttempts: 3,
	// Backoff:     queryBackoff,
	// })

	// Usecases
	userUsecase := userUsecase.NewUserUsecase(s.cfg.App, dbUserRepo)
	projectUsecase := projectUsecase.NewProjectUsecase(s.cfg.App, dbProjectRepo)
	deploymentUsecase := deploymentUsecase.NewDeploymentUsecase(s.cfg.App, dbDeploymentRepo)
	pluginUsecase := pluginUsecase.NewPluginUsecase(s.cfg.App, dbPluginRepo)
	scaffoldRequestUsecase := scaffoldRequestUsecase.NewScaffoldRequestUsecase(s.cfg.App, dbScaffoldRequestRepo)

	// Application middleware
	appMiddleware := middleware.New()

	// Handlers
	userHandler := userHandler.NewUserHandler(s.cfg.App, userUsecase)
	projectHandler := projectHandler.NewProjectHandler(s.cfg.App, projectUsecase)
	deploymentHandler := deploymentHandler.NewDeploymentHandler(s.cfg.App, deploymentUsecase)
	pluginHandler := pluginHandler.NewPluginHandler(s.cfg.App, pluginUsecase)
	scaffoldRequestHandler := scaffoldRequestHandler.NewScaffoldRequestHandler(s.cfg.App, scaffoldRequestUsecase)
	authHandler := authHandler.NewAuthHandler(s.cfg.App, userUsecase)

	return httproute.Dependency{
		Middleware:             appMiddleware,
		UserHandler:            userHandler,
		ProjectHandler:         projectHandler,
		DeploymentHandler:      deploymentHandler,
		PluginHandler:          pluginHandler,
		AuthHandler:            authHandler,
		ScaffoldRequestHandler: scaffoldRequestHandler,
	}, nil
}
