package scaffold

import (
	"context"
	"errors"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

const RunnerName = "scaffold"

type ScaffoldJob struct {
	entity.ScaffoldRequest
}

func (j ScaffoldJob) GetID() uuid.UUID {
	return j.ID
}

type ScaffoldExecutorAdapter struct {
	executor *PythonScaffoldExecutor
}

func NewScaffoldExecutorAdapter(executor *PythonScaffoldExecutor) *ScaffoldExecutorAdapter {
	return &ScaffoldExecutorAdapter{executor: executor}
}

func (a *ScaffoldExecutorAdapter) Execute(ctx context.Context, job *ScaffoldJob) (ScaffoldExecutionResult, error) {
	if job == nil {
		return ScaffoldExecutionResult{}, errors.New("scaffold job is nil")
	}
	return a.executor.Execute(ctx, job)
}

func NewScaffoldPollingRunner(
	observer core.Observability,
	pluginRepository repository.PluginRepository,
	projectRepository repository.ProjectRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
	pollDelay time.Duration,
) (core.Runner, error) {
	executor := NewPythonScaffoldExecutor(pluginRepository, projectRepository)

	// Compose the generic polling runner from queue, state, executor, and observability adapters.
	return core.NewPollingRunner[ScaffoldJob, ScaffoldExecutionResult](
		core.PollingRunnerConfig{
			Name:      RunnerName,
			PollDelay: pollDelay,
		},
		NewScaffoldQueueSourceAdapter(scaffoldRequestRepository),
		NewScaffoldStatePersistence(scaffoldRequestRepository),
		NewScaffoldExecutorAdapter(executor),
		observer,
	)
}
