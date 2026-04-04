package server

import (
	"context"
	httproute "devhub-backend/internal/api/http/route"
	"devhub-backend/internal/config"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/serverutils"
	"fmt"
	"net/http"
	"time"

	infraDB "devhub-backend/internal/infra/db"
	infraLogger "devhub-backend/internal/infra/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg *config.Config
}

func NewServer() *Server {
	return &Server{
		cfg: config.MustConfigure(),
	}
}

func (s *Server) Start() error {
	// Initialize context
	ctx := context.Background()

	// Initialize logger
	logConfig := infraLogger.Config{
		Level:       infraLogger.INFO,
		ServiceName: s.cfg.Service.Name,
		Environment: s.cfg.Service.Env,
	}

	appLogger, err := infraLogger.NewLogger(logConfig)

	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Set default logger config
	if err = infraLogger.SetDefaultLoggerConfig(logConfig); err != nil {
		return fmt.Errorf("failed to set default logger config: %w", err)
	}

	// Initialize database connection
	db, err := infraDB.Connect(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// initialize gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup middlewares
	// middlewares := s.setupMiddlewares(appLogger)

	// Apply middlewares
	// router.Use(middlewares...)

	// Setup not found handler
	router.NoRoute(func(c *gin.Context) {
		httpresponse.Error(c, nil)
	})

	// Setup route dependencies
	deps, err := s.setupRouteDependencies(ctx, appLogger, db)

	if err != nil {
		return fmt.Errorf("failed to setup route dependencies: %w", err)
	}
	// Register application routes
	appRoutes := httproute.NewHTTPRoutes(s.cfg.App, s.cfg.Token, deps)
	appRoutes.RegisterRoutes(router)

	// Create http.Server
	httpServer := &http.Server{
		Addr:              s.cfg.Service.Port,
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
	}

	errCh := make(chan error, 1)

	// Run server in goroutine
	go func() {
		appLogger.Info(ctx, "server started", infraLogger.Fields{
			"service_name": s.cfg.Service.Name,
			"service_env":  s.cfg.Service.Env,
			"service_port": s.cfg.Service.Port,
		})
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	// Wait for shutdown signal then runs provided shutdown tasks in the given order
	shutdownDoneCh := serverutils.GracefulShutdownSystem(
		ctx,
		appLogger,
		errCh,
		30*time.Second,
		[]serverutils.ShutdownTask{
			{
				Name: "HTTP Server",
				Op: func(ctx context.Context) error {
					return httpServer.Shutdown(ctx)
				},
			},
			{
				Name: "Database connection",
				Op: func(ctx context.Context) error {
					return db.Close()
				},
			},
			// Add more shutdown tasks as needed
			// ⚠️ Note: The order of tasks matters.
		},
	)

	// Wait for shutdown to complete
	<-shutdownDoneCh
	appLogger.Info(ctx, "server shutdown complete", nil)
	return nil
}
