package server

import (
	"context"

	"github.com/jmoiron/sqlx"

	authHandler "devhub-backend/internal/api/http/handler/auth"
	deploymentHandler "devhub-backend/internal/api/http/handler/deployment"
	pluginHandler "devhub-backend/internal/api/http/handler/plugin"
	projectHandler "devhub-backend/internal/api/http/handler/project"
	releaseHandler "devhub-backend/internal/api/http/handler/release"
	scaffoldRequestHandler "devhub-backend/internal/api/http/handler/scaffold_request"
	userHandler "devhub-backend/internal/api/http/handler/user"
	"devhub-backend/internal/api/http/middleware"
	httproute "devhub-backend/internal/api/http/route"
	dbDeploymentRepo "devhub-backend/internal/infra/db/repository/deployment"
	dbPluginRepo "devhub-backend/internal/infra/db/repository/plugin"
	dbProjectRepo "devhub-backend/internal/infra/db/repository/project"
	dbRefreshTokenRepo "devhub-backend/internal/infra/db/repository/refresh_token"
	dbReleaseRepo "devhub-backend/internal/infra/db/repository/release"
	dbScaffoldRequestRepo "devhub-backend/internal/infra/db/repository/scaffold_request"
	dbUserRepo "devhub-backend/internal/infra/db/repository/user"
	"devhub-backend/internal/infra/logger"
	authUsecase "devhub-backend/internal/usecase/auth"
	deploymentUsecase "devhub-backend/internal/usecase/deployment"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	projectUsecase "devhub-backend/internal/usecase/project"
	releaseUsecase "devhub-backend/internal/usecase/release"
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
	dbReleaseRepo := dbReleaseRepo.NewReleaseRepository(dbConn)
	dbRefreshTokenRepo := dbRefreshTokenRepo.NewRefreshTokenRepository(dbConn)

	// Query retrier
	// queryBackoff, _ := retry.NewExponentialBackoffStrategy(500*time.Millisecond, 2.0, 5*time.Second)
	// queryRetrier, _ := retry.NewRetrier(retry.Config{
	// MaxAttempts: 3,
	// Backoff:     queryBackoff,
	// })

	// Usecases
	userUsecase := userUsecase.NewUserUsecase(s.cfg.App, dbUserRepo)
	projectUsecase := projectUsecase.NewProjectUsecase(s.cfg.App, dbProjectRepo, dbUserRepo)
	deploymentUsecase := deploymentUsecase.NewDeploymentUsecase(s.cfg.App, dbDeploymentRepo)
	releaseUsecase := releaseUsecase.NewReleaseUsecase(s.cfg.App, dbReleaseRepo)
	pluginUsecase := pluginUsecase.NewPluginUsecase(s.cfg.App, dbPluginRepo)
	scaffoldRequestUsecase := scaffoldRequestUsecase.NewScaffoldRequestUsecase(s.cfg.App, dbScaffoldRequestRepo)
	authUsecase := authUsecase.NewAuthUsecase(s.cfg.Token, dbUserRepo, dbRefreshTokenRepo)

	// Application middleware
	appMiddleware := middleware.New()

	// Handlers
	userHandler := userHandler.NewUserHandler(s.cfg.App, userUsecase)
	projectHandler := projectHandler.NewProjectHandler(s.cfg.App, projectUsecase)
	deploymentHandler := deploymentHandler.NewDeploymentHandler(s.cfg.App, deploymentUsecase)
	releaseHandler := releaseHandler.NewReleaseHandler(s.cfg.App, releaseUsecase)
	pluginHandler := pluginHandler.NewPluginHandler(s.cfg.App, pluginUsecase)
	scaffoldRequestHandler := scaffoldRequestHandler.NewScaffoldRequestHandler(s.cfg.App, scaffoldRequestUsecase)
	authHandler := authHandler.NewAuthHandler(s.cfg.App, authUsecase)

	return httproute.Dependency{
		Middleware:             appMiddleware,
		UserHandler:            userHandler,
		ProjectHandler:         projectHandler,
		DeploymentHandler:      deploymentHandler,
		ReleaseHandler:         releaseHandler,
		PluginHandler:          pluginHandler,
		AuthHandler:            authHandler,
		ScaffoldRequestHandler: scaffoldRequestHandler,
	}, nil
}
