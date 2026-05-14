package deployment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
	"devhub-backend/internal/util/misc"
)

type PythonDeploymentExecutor struct {
	PythonBin        string
	Timeout          time.Duration
	pluginRepository repository.PluginRepository
}

type DeploymentExecutionResult struct {
	ExternalRef string
	CommitSHA   string
}

type deploymentPluginOutput struct {
	Status string `json:"status"`
	Output struct {
		ExternalRef string `json:"external_ref"`
		CommitSHA   string `json:"commit_sha"`
	} `json:"output"`
}

var _ core.Executor[DeploymentJob, DeploymentExecutionResult] = (*DeploymentExecutorAdapter)(nil)

func NewPythonDeploymentExecutor(
	pluginRepository repository.PluginRepository,
) *PythonDeploymentExecutor {
	return &PythonDeploymentExecutor{
		PythonBin:        "python3",
		pluginRepository: pluginRepository,
		Timeout:          10 * time.Minute,
	}
}

func (e *PythonDeploymentExecutor) Execute(
	ctx context.Context,
	job *DeploymentJob,
) (DeploymentExecutionResult, error) {
	if job == nil {
		return DeploymentExecutionResult{}, errors.New("job is nil")
	}

	if e.pluginRepository == nil {
		return DeploymentExecutionResult{}, errors.New("plugin repository is required")
	}

	plugin, err := e.pluginRepository.FindOne(ctx, job.PluginID)
	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) {
			return DeploymentExecutionResult{}, misc.WrapError(
				err,
				errs.NewInternalServerError("failed to find plugin by ID", nil),
			)
		}
		return DeploymentExecutionResult{}, err
	}

	scriptPath := strings.TrimSpace(plugin.Entrypoint)
	if scriptPath == "" {
		return DeploymentExecutionResult{}, errors.New("plugin entrypoint is required")
	}

	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

	in := job.Payload
	stdinBytes, err := json.Marshal(in)
	if err != nil {
		return DeploymentExecutionResult{}, fmt.Errorf("marshal deployment input: %w", err)
	}

	cmd := exec.CommandContext(ctx, e.PythonBin, scriptPath)

	if dir := filepath.Dir(scriptPath); dir != "" && dir != "." {
		cmd.Dir = dir
	}

	cmd.Stdin = bytes.NewReader(stdinBytes)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return DeploymentExecutionResult{}, fmt.Errorf(
			"execute python deployment: %w; stdout=%s; stderr=%s",
			err,
			strings.TrimSpace(stdout.String()),
			strings.TrimSpace(stderr.String()),
		)
	}

	var out deploymentPluginOutput

	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return DeploymentExecutionResult{}, fmt.Errorf(
			"invalid deployment json output: %w; stdout=%s",
			err,
			strings.TrimSpace(stdout.String()),
		)
	}

	if strings.ToLower(strings.TrimSpace(out.Status)) != "ok" {
		return DeploymentExecutionResult{}, errors.New("deployment plugin returned non-ok status")
	}

	return DeploymentExecutionResult{
		ExternalRef: strings.TrimSpace(out.Output.ExternalRef),
		CommitSHA:   strings.TrimSpace(out.Output.CommitSHA),
	}, nil
}
