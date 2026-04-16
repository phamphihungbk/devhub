package worker

import (
	"context"
	"devhub-backend/internal/config"
	"fmt"
	"strings"
	"time"

	"devhub-backend/internal/domain/repository"
	infraLogger "devhub-backend/internal/infra/logger"
	core "devhub-backend/internal/infra/worker/core"
	deployment "devhub-backend/internal/infra/worker/deployment"
	"devhub-backend/internal/infra/worker/release"
	scaffold "devhub-backend/internal/infra/worker/scaffold"
)

// TODO: move to worker configuration
const (
	defaultPollDelay = core.DefaultPollDelay
	RunnerScaffold   = "scaffold"
	RunnerDeployment = "deployment"
	RunnerRelease    = "release"
)

type Dependencies struct {
	cfg                       *config.Config
	logger                    infraLogger.Logger
	pluginRepository          repository.PluginRepository
	projectRepository         repository.ProjectRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
	deploymentRepository      repository.DeploymentRepository
	releaseRepository         repository.ReleaseRepository
}

func NewDependencies(
	cfg *config.Config,
	logger infraLogger.Logger,
	pluginRepository repository.PluginRepository,
	projectRepository repository.ProjectRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
	deploymentRepository repository.DeploymentRepository,
	releaseRepository repository.ReleaseRepository,
) *Dependencies {
	return &Dependencies{
		cfg:                       cfg,
		logger:                    logger,
		pluginRepository:          pluginRepository,
		projectRepository:         projectRepository,
		scaffoldRequestRepository: scaffoldRequestRepository,
		deploymentRepository:      deploymentRepository,
		releaseRepository:         releaseRepository,
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
		RunnerRelease:    buildReleaseRunner,
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

	return scaffold.NewScaffoldPollingRunner(
		observer,
		deps.pluginRepository,
		deps.projectRepository,
		deps.scaffoldRequestRepository,
		cfg.PollDelay,
	)
}

func buildDeploymentRunner(deps *Dependencies, observer Observability, cfg FactoryConfig) (Runner, error) {
	if deps == nil || deps.deploymentRepository == nil {
		return nil, fmt.Errorf("deployment repository is required")
	}

	return deployment.NewDeploymentPollingRunner(
		observer,
		deps.pluginRepository,
		deps.projectRepository,
		deps.deploymentRepository,
		cfg.PollDelay,
	)
}

func buildReleaseRunner(deps *Dependencies, observer Observability, cfg FactoryConfig) (Runner, error) {
	if deps == nil || deps.releaseRepository == nil {
		return nil, fmt.Errorf("release repository is required")
	}

	return release.NewReleasePollingRunner(
		observer,
		deps.pluginRepository,
		deps.projectRepository,
		deps.releaseRepository,
		cfg.PollDelay,
	)
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
