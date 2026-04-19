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
	serviceRepository repository.ServiceRepository
}

type ReleaseExecutionResult struct {
	ExternalRef string
	CommitSHA   string
	FinishedAt  time.Time
}

type releasePluginPayload struct {
	ReleaseID string `json:"release_id"`
	ServiceID string `json:"service_id"`
	PluginID  string `json:"plugin_id"`
	Tag       string `json:"tag"`
	Target    string `json:"target"`
	Name      string `json:"name"`
	Notes     string `json:"notes"`
	RepoURL   string `json:"repo_url"`
}

type releasePluginInput struct {
	Action        string               `json:"action"`
	CorrelationID string               `json:"correlation_id"`
	Payload       releasePluginPayload `json:"payload"`
}

type releasePluginOutput struct {
	Status string `json:"status"`
	Output struct {
		ExternalRef string `json:"external_ref"`
		CommitSHA   string `json:"commit_sha"`
		FinishedAt  string `json:"finished_at"`
	} `json:"output"`
	Error string `json:"error"`
}

var _ core.Executor[ReleaseJob, ReleaseExecutionResult] = (*ReleaseExecutorAdapter)(nil)

func NewPythonReleaseExecutor(
	pluginRepository repository.PluginRepository,
	serviceRepository repository.ServiceRepository,
) *PythonReleaseExecutor {
	return &PythonReleaseExecutor{
		PythonBin:         "python3",
		pluginRepository:  pluginRepository,
		serviceRepository: serviceRepository,
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

	if e.serviceRepository == nil {
		return ReleaseExecutionResult{}, errors.New("service repository is required")
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

	service, err := e.serviceRepository.FindOne(ctx, job.ServiceID)
	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) {
			return ReleaseExecutionResult{}, misc.WrapError(
				err,
				errs.NewInternalServerError("failed to find service by ID", nil),
			)
		}
		return ReleaseExecutionResult{}, err
	}
	if service == nil {
		return ReleaseExecutionResult{}, errors.New("service is required")
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

	payload := releasePluginPayload{
		ReleaseID: job.ID.String(),
		ServiceID: job.ServiceID.String(),
		PluginID:  job.PluginID.String(),
		Tag:       job.Tag,
		Target:    job.Target,
		Name:      job.Name,
		Notes:     job.Notes,
		RepoURL:   service.RepoURL,
	}

	in := releasePluginInput{
		Action:        "release",
		CorrelationID: job.ID.String(),
		Payload:       payload,
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

	var out releasePluginOutput

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
