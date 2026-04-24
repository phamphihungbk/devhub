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
	serviceRepository repository.ServiceRepository
}

type DeploymentExecutionResult struct {
	ExternalRef  string
	CommitSHA    string
	RunnerOutput string
	RunnerError  string
	FinishedAt   time.Time
}

type deploymentPluginPayload struct {
	DeploymentID    string `json:"deployment_id"`
	ProjectID       string `json:"project_id"`
	ServiceID       string `json:"service_id"`
	Service         string `json:"service"`
	PluginID        string `json:"plugin_id"`
	Environment     string `json:"environment"`
	Version         string `json:"version"`
	RepoURL         string `json:"repo_url"`
	SCMAPIURL       string `json:"scm_api_url"`
	SCMToken        string `json:"scm_token"`
	GitopsRepoOwner string `json:"gitops_repo_owner"`
	GitopsRepoName  string `json:"gitops_repo_name"`
	GitopsBranch    string `json:"gitops_branch"`
	GitopsBasePath  string `json:"gitops_base_path"`
	CommitUserName  string `json:"commit_user_name"`
	CommitUserEmail string `json:"commit_user_email"`
	ArgocdServer    string `json:"argocd_server"`
	ArgocdAuthToken string `json:"argocd_auth_token"`
	ArgocdInsecure  bool   `json:"argocd_insecure"`
}

type deploymentPluginInput struct {
	Action        string                  `json:"action"`
	CorrelationID string                  `json:"correlation_id"`
	Payload       deploymentPluginPayload `json:"payload"`
}

type deploymentPluginOutput struct {
	Status string `json:"status"`
	Output struct {
		ExternalRef string `json:"external_ref"`
		CommitSHA   string `json:"commit_sha"`
		FinishedAt  string `json:"finished_at"`
	} `json:"output"`
	Error string `json:"error"`
}

var _ core.Executor[DeploymentJob, DeploymentExecutionResult] = (*DeploymentExecutorAdapter)(nil)

func NewPythonDeploymentExecutor(
	cfg *config.Config,
	pluginRepository repository.PluginRepository,
	serviceRepository repository.ServiceRepository,
) *PythonDeploymentExecutor {
	return &PythonDeploymentExecutor{
		PythonBin:         "python3",
		cfg:               cfg,
		pluginRepository:  pluginRepository,
		serviceRepository: serviceRepository,
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

	if e.serviceRepository == nil {
		return DeploymentExecutionResult{}, errors.New("service repository is required")
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

	service, err := e.serviceRepository.FindOne(ctx, job.ServiceID)
	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) {
			return DeploymentExecutionResult{}, misc.WrapError(
				err,
				errs.NewInternalServerError("failed to find service by ID", nil),
			)
		}
		return DeploymentExecutionResult{}, err
	}
	if service == nil {
		return DeploymentExecutionResult{}, errors.New("service is required")
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

	payload := deploymentPluginPayload{
		DeploymentID:    job.ID.String(),
		ProjectID:       service.ProjectID.String(),
		ServiceID:       job.ServiceID.String(),
		Service:         strings.TrimSpace(service.Name),
		PluginID:        job.PluginID.String(),
		Environment:     job.Environment.String(),
		Version:         job.Version,
		RepoURL:         strings.TrimSpace(service.RepoURL),
		SCMAPIURL:       strings.TrimSpace(e.cfg.ScmConfig.APIURL),
		SCMToken:        strings.TrimSpace(e.cfg.ScmConfig.Token),
		GitopsRepoOwner: strings.TrimSpace(e.cfg.Gitops.RepoOwner),
		GitopsRepoName:  strings.TrimSpace(e.cfg.Gitops.RepoName),
		GitopsBranch:    strings.TrimSpace(e.cfg.Gitops.Branch),
		GitopsBasePath:  strings.TrimSpace(e.cfg.Gitops.BasePath),
		CommitUserName:  strings.TrimSpace(e.cfg.Gitops.CommitUserName),
		CommitUserEmail: strings.TrimSpace(e.cfg.Gitops.CommitUserEmail),
		ArgocdServer:    strings.TrimSpace(e.cfg.ArgoCD.Server),
		ArgocdAuthToken: strings.TrimSpace(e.cfg.ArgoCD.AuthToken),
		ArgocdInsecure:  e.cfg.ArgoCD.Insecure,
	}

	in := deploymentPluginInput{
		Action:        "deploy",
		CorrelationID: job.ID.String(),
		Payload:       payload,
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

	var out deploymentPluginOutput

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
		ExternalRef:  strings.TrimSpace(out.Output.ExternalRef),
		CommitSHA:    strings.TrimSpace(out.Output.CommitSHA),
		RunnerOutput: strings.TrimSpace(stdout.String()),
		RunnerError:  strings.TrimSpace(stderr.String()),
		FinishedAt:   time.Now().UTC(),
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
