package deployment

import (
	"context"
	"devhub-backend/internal/config"
	"errors"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

const RunnerName = "deployment"

type DeploymentJob struct {
	entity.Job
}

func (j DeploymentJob) GetID() uuid.UUID {
	return j.ID
}

type DeploymentExecutorAdapter struct {
	executor *PythonDeploymentExecutor
}

func NewDeploymentExecutorAdapter(executor *PythonDeploymentExecutor) *DeploymentExecutorAdapter {
	return &DeploymentExecutorAdapter{executor: executor}
}

func (a *DeploymentExecutorAdapter) Execute(ctx context.Context, job *DeploymentJob) (DeploymentExecutionResult, error) {
	if job == nil {
		return DeploymentExecutionResult{}, errors.New("deployment job is nil")
	}
	return a.executor.Execute(ctx, job)
}

func NewDeploymentPollingRunner(
	observer core.Observability,
	cfg *config.Config,
	pluginRepository repository.PluginRepository,
	serviceRepository repository.ServiceRepository,
	jobRepository repository.JobRepository,
	deploymentRepository repository.DeploymentRepository,
	releaseRepository repository.ReleaseRepository,
	pollDelay time.Duration,
) (core.Runner, error) {
	_ = cfg
	_ = serviceRepository
	_ = deploymentRepository
	_ = releaseRepository
	executor := NewPythonDeploymentExecutor(pluginRepository)

	return core.NewPollingRunner[DeploymentJob, DeploymentExecutionResult](
		core.PollingRunnerConfig{
			Name:      RunnerName,
			PollDelay: pollDelay,
		},
		NewQueueSourceAdapter(jobRepository),
		NewStatePersistence(jobRepository, deploymentRepository),
		NewDeploymentExecutorAdapter(executor),
		observer,
	)
}
