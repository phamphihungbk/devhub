package scaffold_runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type PythonScaffoldExecutor struct {
	PythonBin  string
	ScriptPath string
	WorkingDir string
	Timeout    time.Duration
}

type ScaffoldExecutionResult struct {
	RepoURL string
}

func NewPythonScaffoldExecutor(scriptPath string) *PythonScaffoldExecutor {
	return &PythonScaffoldExecutor{
		PythonBin:  "python3",
		ScriptPath: scriptPath,
		Timeout:    5 * time.Minute,
	}
}

func (e *PythonScaffoldExecutor) Execute(ctx context.Context, job *ScaffoldJob) (ScaffoldExecutionResult, error) {
	if job == nil {
		return ScaffoldExecutionResult{}, errors.New("job is nil")
	}

	if e.ScriptPath == "" {
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

	if strings.TrimSpace(job.Variables) != "" {
		var vars map[string]any
		if err := json.Unmarshal([]byte(job.Variables), &vars); err == nil {
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

	cmd := exec.CommandContext(ctx, e.PythonBin, e.ScriptPath)

	if e.WorkingDir != "" {
		cmd.Dir = e.WorkingDir
	}

	cmd.Stdin = bytes.NewReader(stdinBytes)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ScaffoldExecutionResult{}, fmt.Errorf("execute python scaffold: %w; stderr=%s", err, strings.TrimSpace(stderr.String()))
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
