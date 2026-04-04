package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	infraLogger "devhub-backend/internal/infra/logger"
	scaffoldRunner "devhub-backend/internal/infra/worker/runner/scaffold"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	defaultPollDelay = 3 * time.Second
	RunnerScaffold   = "scaffold"
	RunnerDeployment = "deployment"
)

type Runner interface {
	Name() string
	Run(ctx context.Context) error
}

type Identifiable interface {
	GetID() uuid.UUID
}

type QueueSourceAdapter[T Identifiable] interface {
	Dequeue(ctx context.Context) (*T, error)
}

type StatePersistence[R any] interface {
	MarkRunning(ctx context.Context, id uuid.UUID) error
	MarkCompleted(ctx context.Context, id uuid.UUID, result R) error
	MarkFailed(ctx context.Context, id uuid.UUID, reason string) error
}

type Executor[T any, R any] interface {
	Execute(ctx context.Context, job *T) (R, error)
}

type BuildRunnersConfig struct {
	WorkerTypes []string
	PollDelay   time.Duration
}

func BuildRunners(logger infraLogger.Logger, db *sqlx.DB) ([]Runner, error) {
	return BuildRunnersWithConfig(logger, db, BuildRunnersConfig{
		WorkerTypes: []string{RunnerScaffold, RunnerDeployment},
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
		case RunnerScaffold:
			runner, err := scaffoldRunner.NewScaffoldPollingRunner(observer, db, cfg.PollDelay)
			if err != nil {
				logger.Warn(context.Background(), "falling back to placeholder scaffold runner", infraLogger.Fields{
					"reason": err.Error(),
				})
				runners = append(runners, newPlaceholderRunner(kind, cfg.PollDelay, logger, db))
				continue
			}
			runners = append(runners, runner)
		case RunnerDeployment:
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
