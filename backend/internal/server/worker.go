package server

import (
	"context"
	"devhub-backend/internal/config"
	"devhub-backend/internal/util/serverutils"
	"errors"
	"fmt"
	"sync"

	"time"

	infraDB "devhub-backend/internal/infra/db"
	dbDeploymentRepo "devhub-backend/internal/infra/db/repository/deployment"
	dbPluginRepo "devhub-backend/internal/infra/db/repository/plugin"
	dbProjectRepo "devhub-backend/internal/infra/db/repository/project"
	dbReleaseRepo "devhub-backend/internal/infra/db/repository/release"
	dbScaffoldRequestRepo "devhub-backend/internal/infra/db/repository/scaffold_request"
	dbServiceRepo "devhub-backend/internal/infra/db/repository/service"
	dbTeamRepo "devhub-backend/internal/infra/db/repository/team"
	infraLogger "devhub-backend/internal/infra/logger"
	infraWorker "devhub-backend/internal/infra/worker"
)

type Worker struct {
	cfg         *config.Config
	concurrency int
	pollDelay   time.Duration
	workerTypes []string
}

type WorkerOption func(*Worker)

func WithWorkerConcurrency(n int) WorkerOption {
	return func(w *Worker) {
		w.concurrency = n
	}
}

func WithWorkerPollInterval(d time.Duration) WorkerOption {
	return func(w *Worker) {
		w.pollDelay = d
	}
}

func WithWorkerTypes(types []string) WorkerOption {
	return func(w *Worker) {
		w.workerTypes = types
	}
}

func NewWorker(opts ...WorkerOption) *Worker {
	w := &Worker{
		cfg:         config.MustConfigure(),
		concurrency: 1,
		pollDelay:   3 * time.Second,
		workerTypes: []string{"scaffold", "deployment"},
	}

	for _, opt := range opts {
		opt(w)
	}

	if w.concurrency <= 0 {
		w.concurrency = 1
	}
	if w.pollDelay <= 0 {
		w.pollDelay = 3 * time.Second
	}

	return w
}

func (w *Worker) Start() error {
	// Initialize context
	ctx := context.Background()

	// Initialize logger
	logConfig := infraLogger.Config{
		Level:       infraLogger.INFO,
		ServiceName: w.cfg.Service.Name,
		Environment: w.cfg.Service.Env,
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
	db, err := infraDB.Connect(w.cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize repository
	dbDeploymentRepo := dbDeploymentRepo.NewDeploymentRepository(db)
	dbPluginRepo := dbPluginRepo.NewPluginRepository(db)
	dbProjectRepo := dbProjectRepo.NewProjectRepository(db)
	dbTeamRepo := dbTeamRepo.NewTeamRepository(db)
	dbScaffoldRequestRepo := dbScaffoldRequestRepo.NewScaffoldRequestRepository(db)
	dbReleaseRepo := dbReleaseRepo.NewReleaseRepository(db)
	dbServiceRepo := dbServiceRepo.NewServiceRepository(db)

	// Initialize runners
	deps := infraWorker.NewDependencies(
		w.cfg,
		appLogger,
		dbPluginRepo,
		dbProjectRepo,
		dbTeamRepo,
		dbScaffoldRequestRepo,
		dbDeploymentRepo,
		dbReleaseRepo,
		dbServiceRepo,
	)

	runners, err := infraWorker.BuildRunnersWithConfig(deps, infraWorker.BuildRunnersConfig{
		WorkerTypes: w.workerTypes,
		PollDelay:   w.pollDelay,
	})

	if len(runners) == 0 {
		db.Close()
		return errors.New("no worker runners configured")
	}

	errCh := make(chan error, 1)
	workerCtx, cancelWorkers := context.WithCancel(ctx)
	var wg sync.WaitGroup

	for _, runner := range runners {
		spawnRunner(workerCtx, appLogger, &wg, errCh, runner, w.concurrency, w.pollDelay)
	}

	shutdownDoneCh := serverutils.GracefulShutdownSystem(
		ctx,
		appLogger,
		errCh,
		30*time.Second,
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
					return db.Close()
				},
			},
			// Add more shutdown tasks as needed
			// ⚠️ Note: The order of tasks matters.
		},
	)

	<-shutdownDoneCh
	appLogger.Info(ctx, "worker shutdown complete", nil)
	return nil
}

func spawnRunner(
	ctx context.Context,
	logger infraLogger.Logger,
	wg *sync.WaitGroup,
	errCh chan<- error,
	runner infraWorker.Runner,
	concurrency int,
	pollDelay time.Duration,
) {
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			logger.Info(ctx, "worker runner started", infraLogger.Fields{
				"runner":      runner.Name(),
				"concurrency": index + 1,
				"poll_delay":  pollDelay.String(),
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
