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
	entity.Deployment
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
	projectRepository repository.ProjectRepository,
	deploymentRepository repository.DeploymentRepository,
	pollDelay time.Duration,
) (core.Runner, error) {
	executor := NewPythonDeploymentExecutor(cfg, pluginRepository, projectRepository)

	return core.NewPollingRunner[DeploymentJob, DeploymentExecutionResult](
		core.PollingRunnerConfig{
			Name:      RunnerName,
			PollDelay: pollDelay,
		},
		NewQueueSourceAdapter(deploymentRepository),
		NewStatePersistence(deploymentRepository),
		NewDeploymentExecutorAdapter(executor),
		observer,
	)
}
