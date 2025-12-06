package server

import (
	"devhub-backend/internal/config"
	"fmt"

	infraDB "devhub-backend/internal/infra/db"
)

type Server struct {
	cfg *config.Config
}

func New() *Server {
	return &Server{
		cfg: config.MustConfigure(),
	}
}

func (s *Server) Start() error {
	// Initialize context
	// ctx := context.Background()

	// // Initialize logger
	// logConfig := logger.Config{
	// 	Level:       logger.INFO,
	// 	ServiceName: s.cfg.Service.Name,
	// 	Environment: s.cfg.Service.Env,
	// }
	// appLogger, err := logger.NewLogger(logConfig)

	// if err != nil {
	// 	return fmt.Errorf("failed to initialize logger: %w", err)
	// }

	// // Set default logger config
	// if err = logger.SetDefaultLoggerConfig(logConfig); err != nil {
	// 	return fmt.Errorf("failed to set default logger config: %w", err)
	// }

	// Initialize database connection
	_, err := infraDB.Connect(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// // initialize gin
	// gin.SetMode(gin.ReleaseMode)
	// router := gin.Default()

	// // Setup middlewares
	// // middlewares := s.setupMiddlewares(appLogger)
	// // Apply middlewares
	// // router.Use(middlewares...)

	// // Setup not found handler
	// router.NoRoute(func(c *gin.Context) {
	// 	httpresponse.Error(c, errsFramework.NewNotFoundError("the requested endpoint is not registered", nil))
	// })

	// // Setup route dependencies
	// deps, err := s.setupRouteDependencies(ctx, tracerProvider, appLogger, db, redisClient)
	// if err != nil {
	// 	return fmt.Errorf("failed to setup route dependencies: %w", err)
	// }
	// Register application routes
	// appRoutes := httproute.NewHTTPRoutes(s.cfg.App, deps)
	// appRoutes.RegisterRoutes(router)

	// Create http.Server
	// httpServer := &http.Server{
	// 	Addr:              s.cfg.Service.Port,
	// 	Handler:           router,
	// 	ReadHeaderTimeout: 15 * time.Second,
	// }

	// errCh := make(chan error, 1)

	// // Run server in goroutine
	// go func() {
	// 	appLogger.Info(ctx, "server started", logger.Fields{
	// 		"service_name": s.cfg.Service.Name,
	// 		"service_env":  s.cfg.Service.Env,
	// 		"service_port": s.cfg.Service.Port,
	// 	})
	// 	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		errCh <- fmt.Errorf("http server error: %w", err)
	// 	}
	// }()

	// // Wait for shutdown signal then runs provided shutdown tasks in the given order
	// shutdownDoneCh := serverutils.GracefulShutdownSystem(
	// 	ctx,
	// 	appLogger,
	// 	errCh,
	// 	30*time.Second,
	// 	[]serverutils.ShutdownTask{
	// 		{
	// 			Name: "HTTP Server",
	// 			Op: func(ctx context.Context) error {
	// 				return httpServer.Shutdown(ctx)
	// 			},
	// 		},
	// 		{
	// 			Name: "Tracer Provider",
	// 			Op: func(ctx context.Context) error {
	// 				return tracerProvider.Shutdown(ctx)
	// 			},
	// 		},
	// 		{
	// 			Name: "Redis client",
	// 			Op: func(ctx context.Context) error {
	// 				return redisClient.Close()
	// 			},
	// 		},
	// 		{
	// 			Name: "Database connection",
	// 			Op: func(ctx context.Context) error {
	// 				return db.Close()
	// 			},
	// 		},
	// 		// Add more shutdown tasks as needed
	// 		// ⚠️ Note: The order of tasks matters.
	// 	},
	// )

	// // Wait for shutdown to complete
	// <-shutdownDoneCh
	// appLogger.Info(ctx, "server shutdown complete", nil)
	return nil
}
