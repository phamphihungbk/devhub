package scaffold_runner

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	infraWorker "devhub-backend/internal/infra/worker"

	"github.com/jmoiron/sqlx"
)

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

func NewScaffoldPollingRunner(observer infraWorker.Observability, db *sqlx.DB, pollDelay time.Duration) (infraWorker.Runner, error) {
	// TODO: resolve the script path from plugin configuration in the database.
	scriptPath := strings.TrimSpace(os.Getenv("WORKER_SCAFFOLD_SCRIPT"))

	if scriptPath == "" {
		return nil, errors.New("WORKER_SCAFFOLD_SCRIPT is required for scaffold runner")
	}

	executor := NewPythonScaffoldExecutor(scriptPath)

	if workDir := strings.TrimSpace(os.Getenv("WORKER_SCAFFOLD_WORKDIR")); workDir != "" {
		executor.WorkingDir = workDir
	}

	// Compose the generic polling runner from queue, state, executor, and observability adapters.
	return infraWorker.NewPollingRunner[ScaffoldJob, ScaffoldExecutionResult](
		infraWorker.PollingRunnerConfig{
			Name:      infraWorker.runnerScaffold,
			PollDelay: pollDelay,
		},
		NewScaffoldQueueSourceAdapter(db),
		NewScaffoldStatePersistence(db),
		NewScaffoldExecutorAdapter(executor),
		observer,
	)
}
