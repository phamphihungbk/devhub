package scaffold

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

type PythonScaffoldExecutor struct {
	PythonBin        string
	Timeout          time.Duration
	pluginRepository repository.PluginRepository
}

type ScaffoldExecutionResult struct {
	RepoURL string
}

var _ core.Executor[ScaffoldJob, ScaffoldExecutionResult] = (*ScaffoldExecutorAdapter)(nil)

func NewPythonScaffoldExecutor(pluginRepository repository.PluginRepository) *PythonScaffoldExecutor {
	return &PythonScaffoldExecutor{
		PythonBin:        "python3",
		pluginRepository: pluginRepository,
		Timeout:          5 * time.Minute,
	}
}

func (e *PythonScaffoldExecutor) Execute(ctx context.Context, job *ScaffoldJob) (ScaffoldExecutionResult, error) {
	if job == nil {
		return ScaffoldExecutionResult{}, errors.New("job is nil")
	}

	if e.pluginRepository == nil {
		return ScaffoldExecutionResult{}, errors.New("plugin repository is required")
	}

	plugin, err := e.pluginRepository.FindOne(ctx, job.PluginID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return ScaffoldExecutionResult{}, misc.WrapError(err, errs.NewInternalServerError("failed to find plugin by ID", nil))
		}
		return ScaffoldExecutionResult{}, err // Return the NotFoundError directly
	}

	scriptPath := strings.TrimSpace(plugin.Entrypoint)

	if scriptPath == "" {
		return ScaffoldExecutionResult{}, errors.New("script path is required")
	}

	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

	payload := map[string]any{
		"scaffold_request_id": job.ID.String(),
		"project_id":          job.ProjectID.String(),
		"template":            job.Template,
		"environment":         job.Environment,
	}

	if job.Variables.String() != "" {
		var vars map[string]any
		if err := json.Unmarshal([]byte(job.Variables.String()), &vars); err == nil {
			for k, v := range vars {
				payload[k] = v
			}
		}
	}

	in := map[string]any{
		"action":         "scaffold",
		"correlation_id": job.ID.String(),
		"payload":        payload,
	}

	stdinBytes, err := json.Marshal(in)

	if err != nil {
		return ScaffoldExecutionResult{}, fmt.Errorf("marshal scaffold input: %w", err)
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
		return ScaffoldExecutionResult{}, fmt.Errorf(
			"execute python scaffold: %w; stdout=%s; stderr=%s",
			err,
			strings.TrimSpace(stdout.String()),
			strings.TrimSpace(stderr.String()),
		)
	}

	var out struct {
		Status string `json:"status"`
		Output struct {
			RepoURL string `json:"repo_url"`
			Path    string `json:"path"`
		} `json:"output"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return ScaffoldExecutionResult{}, fmt.Errorf("invalid scaffold json output: %w", err)
	}

	if strings.ToLower(out.Status) != "ok" {
		return ScaffoldExecutionResult{}, fmt.Errorf("plugin returned non-ok status")
	}

	repoURL := strings.TrimSpace(out.Output.RepoURL)

	if repoURL == "" {
		repoURL = strings.TrimSpace(out.Output.Path)
	}

	if repoURL == "" {
		return ScaffoldExecutionResult{}, errors.New("plugin output missing repo_url/path")
	}

	return ScaffoldExecutionResult{RepoURL: repoURL}, nil
}
