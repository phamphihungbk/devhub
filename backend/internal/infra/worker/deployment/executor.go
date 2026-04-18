package deployment

import (
	"bytes"
	"context"
	"devhub-backend/internal/config"
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
	PythonBin         string
	Timeout           time.Duration
	cfg               *config.Config
	pluginRepository  repository.PluginRepository
	projectRepository repository.ProjectRepository
}

type DeploymentExecutionResult struct {
	ExternalRef string
	CommitSHA   string
	FinishedAt  time.Time
}

var _ core.Executor[DeploymentJob, DeploymentExecutionResult] = (*DeploymentExecutorAdapter)(nil)

func NewPythonDeploymentExecutor(
	cfg *config.Config,
	pluginRepository repository.PluginRepository,
	projectRepository repository.ProjectRepository,
) *PythonDeploymentExecutor {
	return &PythonDeploymentExecutor{
		PythonBin:         "python3",
		cfg:               cfg,
		pluginRepository:  pluginRepository,
		projectRepository: projectRepository,
		Timeout:           10 * time.Minute,
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

	if e.cfg == nil {
		return DeploymentExecutionResult{}, errors.New("config is required")
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

	payload := map[string]any{
		"deployment_id":     job.ID.String(),
		"service_id":        job.ServiceID.String(),
		"plugin_id":         job.PluginID.String(),
		"environment":       job.Environment,
		"version":           job.Version,
		"scm_api_url":       strings.TrimSpace(e.cfg.ScmConfig.APIURL),
		"scm_token":         strings.TrimSpace(e.cfg.ScmConfig.Token),
		"gitops_repo_owner": strings.TrimSpace(e.cfg.Gitops.RepoOwner),
		"gitops_repo_name":  strings.TrimSpace(e.cfg.Gitops.RepoName),
		"gitops_branch":     strings.TrimSpace(e.cfg.Gitops.Branch),
		"gitops_base_path":  strings.TrimSpace(e.cfg.Gitops.BasePath),
		"commit_user_name":  strings.TrimSpace(e.cfg.Gitops.CommitUserName),
		"commit_user_email": strings.TrimSpace(e.cfg.Gitops.CommitUserEmail),
		"argocd_server":     strings.TrimSpace(e.cfg.ArgoCD.Server),
		"argocd_auth_token": strings.TrimSpace(e.cfg.ArgoCD.AuthToken),
		"argocd_insecure":   e.cfg.ArgoCD.Insecure,
	}

	in := map[string]any{
		"action":         "deploy",
		"correlation_id": job.ID.String(),
		"payload":        payload,
	}

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
		return DeploymentExecutionResult{}, fmt.Errorf(
			"invalid deployment json output: %w; stdout=%s",
			err,
			strings.TrimSpace(stdout.String()),
		)
	}

	if strings.ToLower(strings.TrimSpace(out.Status)) != "ok" {
		if strings.TrimSpace(out.Error) != "" {
			return DeploymentExecutionResult{}, fmt.Errorf("deployment plugin failed: %s", out.Error)
		}
		return DeploymentExecutionResult{}, errors.New("deployment plugin returned non-ok status")
	}

	result := DeploymentExecutionResult{
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
		result.ExternalRef = fmt.Sprintf("%s-%s", job.ServiceID.String(), job.Environment)
	}

	return result, nil
}
