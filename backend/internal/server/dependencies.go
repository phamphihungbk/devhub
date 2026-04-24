package server

import (
	"context"

	"github.com/jmoiron/sqlx"

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
	httproute "devhub-backend/internal/api/http/route"
	dbApprovalRepo "devhub-backend/internal/infra/db/repository/approval"
	dbDeploymentRepo "devhub-backend/internal/infra/db/repository/deployment"
	dbPluginRepo "devhub-backend/internal/infra/db/repository/plugin"
	dbProjectRepo "devhub-backend/internal/infra/db/repository/project"
	dbRefreshTokenRepo "devhub-backend/internal/infra/db/repository/refresh_token"
	dbReleaseRepo "devhub-backend/internal/infra/db/repository/release"
	dbScaffoldRequestRepo "devhub-backend/internal/infra/db/repository/scaffold_request"
	dbServiceRepo "devhub-backend/internal/infra/db/repository/service"
	dbTeamRepo "devhub-backend/internal/infra/db/repository/team"
	dbUserRepo "devhub-backend/internal/infra/db/repository/user"
	"devhub-backend/internal/infra/logger"
	approvalUsecase "devhub-backend/internal/usecase/approval"
	authUsecase "devhub-backend/internal/usecase/auth"
	deploymentUsecase "devhub-backend/internal/usecase/deployment"
	pluginUsecase "devhub-backend/internal/usecase/plugin"
	projectUsecase "devhub-backend/internal/usecase/project"
	releaseUsecase "devhub-backend/internal/usecase/release"
	scaffoldRequestUsecase "devhub-backend/internal/usecase/scaffold_request"
	serviceUsecase "devhub-backend/internal/usecase/service"
	teamUsecase "devhub-backend/internal/usecase/team"
	userUsecase "devhub-backend/internal/usecase/user"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

//nolint:unparam
func (s *Server) setupRouteDependencies(ctx context.Context, tracerProvider *sdktrace.TracerProvider, appLogger logger.Logger, dbConn *sqlx.DB) (httproute.Dependency, error) {
	// Transactor factory
	// transactorFactory := infraDB.NewSqlxTransactorFactory(dbConn)

	// DB Repositories
	dbUserRepo := dbUserRepo.NewUserRepository(dbConn)
	dbApprovalRepo := dbApprovalRepo.NewApprovalRepository(dbConn)
	dbProjectRepo := dbProjectRepo.NewProjectRepository(dbConn)
	dbDeploymentRepo := dbDeploymentRepo.NewDeploymentRepository(dbConn)
	dbPluginRepo := dbPluginRepo.NewPluginRepository(dbConn)
	dbScaffoldRequestRepo := dbScaffoldRequestRepo.NewScaffoldRequestRepository(dbConn)
	dbReleaseRepo := dbReleaseRepo.NewReleaseRepository(dbConn)
	dbServiceRepo := dbServiceRepo.NewServiceRepository(dbConn)
	dbTeamRepo := dbTeamRepo.NewTeamRepository(dbConn)
	dbRefreshTokenRepo := dbRefreshTokenRepo.NewRefreshTokenRepository(dbConn)

	// Usecases
	approvalUsecase := approvalUsecase.NewApprovalUsecase(s.cfg.App, dbApprovalRepo, dbScaffoldRequestRepo)
	userUsecase := userUsecase.NewUserUsecase(s.cfg.App, dbUserRepo)
	projectUsecase := projectUsecase.NewProjectUsecase(s.cfg.App, dbProjectRepo, dbUserRepo)
	deploymentUsecase := deploymentUsecase.NewDeploymentUsecase(s.cfg.App, dbApprovalRepo, dbDeploymentRepo)
	releaseUsecase := releaseUsecase.NewReleaseUsecase(s.cfg.App, dbReleaseRepo)
	pluginUsecase := pluginUsecase.NewPluginUsecase(s.cfg.App, dbPluginRepo)
	scaffoldRequestUsecase := scaffoldRequestUsecase.NewScaffoldRequestUsecase(s.cfg.App, dbApprovalRepo, dbScaffoldRequestRepo)
	serviceUsecase := serviceUsecase.NewServiceUsecase(s.cfg.App, dbServiceRepo)
	teamUsecase := teamUsecase.NewTeamUsecase(s.cfg.App, dbTeamRepo)
	authUsecase := authUsecase.NewAuthUsecase(s.cfg.Token, dbUserRepo, dbRefreshTokenRepo)

	// Application middleware
	appMiddleware := middleware.New()

	// Handlers
	approvalHandler := approvalHandler.NewApprovalHandler(s.cfg.App, approvalUsecase)
	userHandler := userHandler.NewUserHandler(s.cfg.App, userUsecase)
	projectHandler := projectHandler.NewProjectHandler(s.cfg.App, projectUsecase)
	deploymentHandler := deploymentHandler.NewDeploymentHandler(s.cfg.App, deploymentUsecase)
	releaseHandler := releaseHandler.NewReleaseHandler(s.cfg.App, releaseUsecase)
	pluginHandler := pluginHandler.NewPluginHandler(s.cfg.App, pluginUsecase)
	scaffoldRequestHandler := scaffoldRequestHandler.NewScaffoldRequestHandler(s.cfg.App, scaffoldRequestUsecase)
	serviceHandler := serviceHandler.NewServiceHandler(s.cfg.App, serviceUsecase)
	teamHandler := teamHandler.NewTeamHandler(s.cfg.App, teamUsecase)
	authHandler := authHandler.NewAuthHandler(s.cfg.App, authUsecase)

	return httproute.Dependency{
		Middleware:             appMiddleware,
		ApprovalHandler:        approvalHandler,
		UserHandler:            userHandler,
		ProjectHandler:         projectHandler,
		DeploymentHandler:      deploymentHandler,
		ReleaseHandler:         releaseHandler,
		PluginHandler:          pluginHandler,
		AuthHandler:            authHandler,
		ScaffoldRequestHandler: scaffoldRequestHandler,
		ServiceHandler:         serviceHandler,
		TeamHandler:            teamHandler,
	}, nil
}
