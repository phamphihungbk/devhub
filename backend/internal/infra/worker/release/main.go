package release

import (
	"context"
	"errors"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

const RunnerName = "release"

type ReleaseJob struct {
	entity.Release
}

func (j ReleaseJob) GetID() uuid.UUID {
	return j.ID
}

type ReleaseExecutorAdapter struct {
	executor *PythonReleaseExecutor
}

func NewReleaseExecutorAdapter(executor *PythonReleaseExecutor) *ReleaseExecutorAdapter {
	return &ReleaseExecutorAdapter{executor: executor}
}

func (a *ReleaseExecutorAdapter) Execute(ctx context.Context, job *ReleaseJob) (ReleaseExecutionResult, error) {
	if job == nil {
		return ReleaseExecutionResult{}, errors.New("release job is nil")
	}
	return a.executor.Execute(ctx, job)
}

func NewReleasePollingRunner(
	observer core.Observability,
	pluginRepository repository.PluginRepository,
	projectRepository repository.ProjectRepository,
	releaseRepository repository.ReleaseRepository,
	pollDelay time.Duration,
) (core.Runner, error) {
	executor := NewPythonReleaseExecutor(pluginRepository, projectRepository)

	return core.NewPollingRunner[ReleaseJob, ReleaseExecutionResult](
		core.PollingRunnerConfig{
			Name:      RunnerName,
			PollDelay: pollDelay,
		},
		NewQueueSourceAdapter(releaseRepository),
		NewStatePersistence(releaseRepository),
		NewReleaseExecutorAdapter(executor),
		observer,
	)
}
