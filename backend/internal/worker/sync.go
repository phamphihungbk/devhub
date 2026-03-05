package worker

import (
	"context"
	"fmt"
	"time"

	"devhub-backend/internal/config"
	infraDB "devhub-backend/internal/infra/db"
	infraLogger "devhub-backend/internal/infra/logger"
	infraWorkerQueue "devhub-backend/internal/infra/worker/queue"
	infraWorkerRuntime "devhub-backend/internal/infra/worker/runtime"
	"devhub-backend/internal/util/serverutils"
)

type SyncConfig struct {
	PollInterval  time.Duration
	PluginTimeout time.Duration
	PluginCommand []string
	MaxWorkers    int
	MaxAttempts   int
}

type SyncWorker struct {
	cfg        *config.Config
	syncConfig SyncConfig
}

func NewSync(syncConfig SyncConfig) *SyncWorker {
	return &SyncWorker{
		cfg:        config.MustConfigure(),
		syncConfig: syncConfig,
	}
}

func (w *SyncWorker) Start() error {
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

	if err = infraLogger.SetDefaultLoggerConfig(logConfig); err != nil {
		return fmt.Errorf("failed to set default logger config: %w", err)
	}

	db, err := infraDB.Connect(w.cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	appLogger.Info(ctx, "worker sync initialized", infraLogger.Fields{
		"poll_interval":  w.syncConfig.PollInterval.String(),
		"plugin_timeout": w.syncConfig.PluginTimeout.String(),
		"plugin_command": w.syncConfig.PluginCommand,
		"max_workers":    w.syncConfig.MaxWorkers,
		"max_attempts":   w.syncConfig.MaxAttempts,
		"worker_mode":    "sync",
	})

	runtimeWorker := infraWorkerRuntime.New(
		infraWorkerRuntime.Config{
			PollInterval: w.syncConfig.PollInterval,
			MaxWorkers:   w.syncConfig.MaxWorkers,
			RetryPolicy: infraWorkerRuntime.RetryPolicy{
				MaxAttempts: w.syncConfig.MaxAttempts,
			},
		},
		emptyQueue{},
		noopRunner{},
		appLogger,
	)

	workerCtx, cancelWorker := context.WithCancel(ctx)
	workerDone := make(chan struct{})
	errCh := make(chan error, 1)

	go func() {
		defer close(workerDone)
		if err := runtimeWorker.Start(workerCtx); err != nil {
			errCh <- err
		}
	}()

	shutdownDoneCh := serverutils.GracefulShutdownSystem(
		ctx,
		appLogger,
		errCh,
		30*time.Second,
		[]serverutils.ShutdownTask{
			{
				Name: "Worker runtime",
				Op: func(ctx context.Context) error {
					cancelWorker()
					select {
					case <-workerDone:
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				},
			},
			{
				Name: "Database connection",
				Op: func(ctx context.Context) error {
					return db.Close()
				},
			},
		},
	)

	<-shutdownDoneCh
	appLogger.Info(ctx, "worker shutdown complete", nil)
	return nil
}
