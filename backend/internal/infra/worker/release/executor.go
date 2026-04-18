package release

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

type PythonReleaseExecutor struct {
	PythonBin         string
	Timeout           time.Duration
	pluginRepository  repository.PluginRepository
}

type ReleaseExecutionResult struct {
	ExternalRef string
	CommitSHA   string
	FinishedAt  time.Time
}

var _ core.Executor[ReleaseJob, ReleaseExecutionResult] = (*ReleaseExecutorAdapter)(nil)

func NewPythonReleaseExecutor(
	pluginRepository repository.PluginRepository,
) *PythonReleaseExecutor {
	return &PythonReleaseExecutor{
		PythonBin:         "python3",
		pluginRepository:  pluginRepository,
		Timeout:           10 * time.Minute,
	}
}

func (e *PythonReleaseExecutor) Execute(
	ctx context.Context,
	job *ReleaseJob,
) (ReleaseExecutionResult, error) {
	if job == nil {
		return ReleaseExecutionResult{}, errors.New("job is nil")
	}

	if e.pluginRepository == nil {
		return ReleaseExecutionResult{}, errors.New("plugin repository is required")
	}

	plugin, err := e.pluginRepository.FindOne(ctx, job.PluginID)
	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) {
			return ReleaseExecutionResult{}, misc.WrapError(
				err,
				errs.NewInternalServerError("failed to find plugin by ID", nil),
			)
		}
		return ReleaseExecutionResult{}, err
	}

	scriptPath := strings.TrimSpace(plugin.Entrypoint)
	if scriptPath == "" {
		return ReleaseExecutionResult{}, errors.New("plugin entrypoint is required")
	}

	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

	payload := map[string]any{
		"release_id":    job.ID.String(),
		"service_id":    job.ServiceID.String(),
		"plugin_id":     job.PluginID.String(),
		"tag":           job.Tag,
		"target":        job.Target,
		"name":          job.Name,
		"notes":         job.Notes,
	}

	in := map[string]any{
		"action":         "release",
		"correlation_id": job.ID.String(),
		"payload":        payload,
	}

	stdinBytes, err := json.Marshal(in)
	if err != nil {
		return ReleaseExecutionResult{}, fmt.Errorf("marshal release input: %w", err)
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
		return ReleaseExecutionResult{}, fmt.Errorf(
			"execute python release: %w; stdout=%s; stderr=%s",
			err,
			strings.TrimSpace(stdout.String()),
			strings.TrimSpace(stderr.String()),
		)
	}

	var out struct {
		Status string `json:"status"`
		Output struct {
			ExternalRef string `json:"external_ref"`
			CommitSHA   string `json:"commit_sha"`
			FinishedAt  string `json:"finished_at"`
		} `json:"output"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return ReleaseExecutionResult{}, fmt.Errorf(
			"invalid release json output: %w; stdout=%s",
			err,
			strings.TrimSpace(stdout.String()),
		)
	}

	if strings.ToLower(strings.TrimSpace(out.Status)) != "ok" {
		if strings.TrimSpace(out.Error) != "" {
			return ReleaseExecutionResult{}, fmt.Errorf("release plugin failed: %s", out.Error)
		}
		return ReleaseExecutionResult{}, errors.New("release plugin returned non-ok status")
	}

	result := ReleaseExecutionResult{
		ExternalRef: strings.TrimSpace(out.Output.ExternalRef),
		CommitSHA:   strings.TrimSpace(out.Output.CommitSHA),
		FinishedAt:  time.Now().UTC(),
	}

	if ts := strings.TrimSpace(out.Output.FinishedAt); ts != "" {
		if parsed, err := time.Parse(time.RFC3339, ts); err == nil {
			result.FinishedAt = parsed.UTC()
		}
	}

	if result.ExternalRef == "" {
		result.ExternalRef = job.Tag
	}

	return result, nil
}
