package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"devhub-backend/internal/domain/repository"
	infraLogger "devhub-backend/internal/infra/logger"
	core "devhub-backend/internal/infra/worker/core"
	scaffold "devhub-backend/internal/infra/worker/scaffold"
)

const (
	defaultPollDelay = core.DefaultPollDelay
	RunnerScaffold   = "scaffold"
	RunnerDeployment = "deployment"
)

type Dependencies struct {
	logger                    infraLogger.Logger
	pluginRepository          repository.PluginRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

func NewDependencies(
	logger infraLogger.Logger,
	pluginRepository repository.PluginRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
) *Dependencies {
	return &Dependencies{
		logger:                    logger,
		pluginRepository:          pluginRepository,
		scaffoldRequestRepository: scaffoldRequestRepository,
	}
}

type Runner = core.Runner

type BuildRunnersConfig struct {
	WorkerTypes []string
	PollDelay   time.Duration
}

type FactoryConfig struct {
	PollDelay time.Duration
}

type RunnerFactory func(deps *Dependencies, observer Observability, cfg FactoryConfig) (Runner, error)

func BuildRunnersWithConfig(deps *Dependencies, cfg BuildRunnersConfig) ([]Runner, error) {
	if deps == nil {
		return nil, fmt.Errorf("worker dependencies are required")
	}
	if deps.logger == nil {
		return nil, fmt.Errorf("worker logger dependency is required")
	}
	if cfg.PollDelay <= 0 {
		cfg.PollDelay = defaultPollDelay
	}

	workerTypes := normalizeWorkerTypes(cfg.WorkerTypes)

	if len(workerTypes) == 0 {
		return []Runner{}, nil
	}

	runners := make([]Runner, 0, len(workerTypes))
	observer := NewLoggerObservability(deps.logger)

	factories := map[string]RunnerFactory{
		RunnerScaffold:   buildScaffoldRunner,
		RunnerDeployment: buildDeploymentRunner,
	}

	for _, kind := range workerTypes {
		factory, ok := factories[kind]
		if !ok {
			return nil, fmt.Errorf("unsupported worker type: %s", kind)
		}

		runner, err := factory(deps, observer, FactoryConfig{PollDelay: cfg.PollDelay})
		if err != nil {
			deps.logger.Warn(context.Background(), "falling back to placeholder runner", infraLogger.Fields{
				"runner": kind,
				"reason": err.Error(),
			})
			runners = append(runners, newPlaceholderRunner(kind, cfg.PollDelay, deps.logger))
			continue
		}

		runners = append(runners, runner)
	}

	return runners, nil
}

func buildScaffoldRunner(deps *Dependencies, observer Observability, cfg FactoryConfig) (Runner, error) {
	if deps == nil {
		return nil, fmt.Errorf("worker dependencies are required")
	}

	return scaffold.NewScaffoldPollingRunner(observer, deps.pluginRepository, deps.scaffoldRequestRepository, cfg.PollDelay)
}

func buildDeploymentRunner(deps *Dependencies, _ Observability, cfg FactoryConfig) (Runner, error) {
	if deps == nil || deps.logger == nil {
		return nil, fmt.Errorf("worker logger dependency is required")
	}

	return newPlaceholderRunner(RunnerDeployment, cfg.PollDelay, deps.logger), nil
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
