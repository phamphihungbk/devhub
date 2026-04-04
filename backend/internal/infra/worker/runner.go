package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	infraLogger "devhub-backend/internal/infra/logger"
	scaffoldRunner "devhub-backend/internal/infra/worker/runner/scaffold"

	"github.com/jmoiron/sqlx"
)

const (
	defaultPollDelay = 3 * time.Second
	runnerScaffold   = "scaffold"
	runnerDeployment = "deployment"
)

type Runner interface {
	Name() string
	Run(ctx context.Context) error
}

type BuildRunnersConfig struct {
	WorkerTypes []string
	PollDelay   time.Duration
}

func BuildRunners(logger infraLogger.Logger, db *sqlx.DB) ([]Runner, error) {
	return BuildRunnersWithConfig(logger, db, BuildRunnersConfig{
		WorkerTypes: []string{runnerScaffold, runnerDeployment},
		PollDelay:   defaultPollDelay,
	})
}

func BuildRunnersWithConfig(logger infraLogger.Logger, db *sqlx.DB, cfg BuildRunnersConfig) ([]Runner, error) {
	if cfg.PollDelay <= 0 {
		cfg.PollDelay = defaultPollDelay
	}

	workerTypes := normalizeWorkerTypes(cfg.WorkerTypes)

	if len(workerTypes) == 0 {
		return []Runner{}, nil
	}

	runners := make([]Runner, 0, len(workerTypes))
	observer := NewLoggerObservability(logger)

	for _, kind := range workerTypes {
		switch kind {
		case runnerScaffold:
			runner, err := scaffoldRunner.NewScaffoldPollingRunner(observer, db, cfg.PollDelay)
			if err != nil {
				logger.Warn(context.Background(), "falling back to placeholder scaffold runner", infraLogger.Fields{
					"reason": err.Error(),
				})
				runners = append(runners, newPlaceholderRunner(kind, cfg.PollDelay, logger, db))
				continue
			}
			runners = append(runners, runner)
		case runnerDeployment:
			runners = append(runners, newPlaceholderRunner(kind, cfg.PollDelay, logger, db))
		default:
			return nil, fmt.Errorf("unsupported worker type: %s", kind)
		}
	}

	return runners, nil
}

func normalizeWorkerTypes(types []string) []string {
	seen := make(map[string]struct{}, len(types))
	result := make([]string, 0, len(types))

	for _, t := range types {
		v := strings.ToLower(strings.TrimSpace(t))
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}

	return result
}

type placeholderRunner struct {
	name      string
	pollDelay time.Duration
	logger    infraLogger.Logger
	db        *sqlx.DB
}

func newPlaceholderRunner(name string, pollDelay time.Duration, logger infraLogger.Logger, db *sqlx.DB) *placeholderRunner {
	return &placeholderRunner{
		name:      name,
		pollDelay: pollDelay,
		logger:    logger,
		db:        db,
	}
}

func (r *placeholderRunner) Name() string {
	return r.name
}

func (r *placeholderRunner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.pollDelay)
	defer ticker.Stop()

	r.logger.Warn(ctx, "runner is in placeholder mode; wire queue/executor to enable processing", infraLogger.Fields{
		"runner": r.name,
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_ = r.db
		}
	}
}
