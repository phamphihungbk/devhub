package worker

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/util/serverutils"
	"errors"
	"fmt"
	"sync"
	"time"

	infraDB "devhub-backend/internal/infra/db"
	infraLogger "devhub-backend/internal/infra/logger"

	"github.com/jmoiron/sqlx"
)

const defaultShutdownTimeout = 30 * time.Second

type Runtime struct {
	logger infraLogger.Logger
	db     *sqlx.DB
}

type RunConfig struct {
	Concurrency     int
	PollDelay       time.Duration
	ShutdownTimeout time.Duration
}

func Bootstrap(cfg *config.Config) (*Runtime, error) {
	logConfig := infraLogger.Config{
		Level:       infraLogger.INFO,
		ServiceName: cfg.Service.Name,
		Environment: cfg.Service.Env,
	}

	appLogger, err := infraLogger.NewLogger(logConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err = infraLogger.SetDefaultLoggerConfig(logConfig); err != nil {
		return nil, fmt.Errorf("failed to set default logger config: %w", err)
	}

	db, err := infraDB.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Runtime{
		logger: appLogger,
		db:     db,
	}, nil
}

func (r *Runtime) Logger() infraLogger.Logger {
	return r.logger
}

func (r *Runtime) DB() *sqlx.DB {
	return r.db
}

func Run(ctx context.Context, rt *Runtime, runners []Runner, cfg RunConfig) error {
	if rt == nil {
		return errors.New("runtime is required")
	}
	if len(runners) == 0 {
		_ = rt.db.Close()
		return errors.New("no worker runners configured")
	}

	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 1
	}

	if cfg.PollDelay <= 0 {
		cfg.PollDelay = defaultPollDelay
	}

	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = defaultShutdownTimeout
	}

	errCh := make(chan error, 1)
	workerCtx, cancelWorkers := context.WithCancel(ctx)
	var wg sync.WaitGroup

	for _, runner := range runners {
		spawnRunner(workerCtx, rt.logger, &wg, errCh, runner, cfg)
	}

	shutdownDoneCh := serverutils.GracefulShutdownSystem(
		ctx,
		rt.logger,
		errCh,
		cfg.ShutdownTimeout,
		[]serverutils.ShutdownTask{
			{
				Name: "Worker runners",
				Op: func(ctx context.Context) error {
					cancelWorkers()
					return waitGroupWithContext(ctx, &wg)
				},
			},
			{
				Name: "Database connection",
				Op: func(ctx context.Context) error {
					return rt.db.Close()
				},
			},
		},
	)

	<-shutdownDoneCh
	rt.logger.Info(ctx, "worker shutdown complete", nil)
	return nil
}

func spawnRunner(
	ctx context.Context,
	logger infraLogger.Logger,
	wg *sync.WaitGroup,
	errCh chan<- error,
	runner Runner,
	cfg RunConfig,
) {
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			logger.Info(ctx, "worker runner started", infraLogger.Fields{
				"runner":      runner.Name(),
				"concurrency": index + 1,
				"poll_delay":  cfg.PollDelay.String(),
			})
			if runErr := runner.Run(ctx); runErr != nil && !errors.Is(runErr, context.Canceled) {
				select {
				case errCh <- fmt.Errorf("runner %s failed: %w", runner.Name(), runErr):
				default:
				}
			}
		}(i)
	}
}

func waitGroupWithContext(ctx context.Context, wg *sync.WaitGroup) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
