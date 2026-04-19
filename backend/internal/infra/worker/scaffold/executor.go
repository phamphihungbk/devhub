package scaffold

import (
	"bytes"
	"context"
	"devhub-backend/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
	"devhub-backend/internal/util/misc"

	"github.com/google/uuid"
)

const DEFAULT_IMAGE_TAG = "latest"

type PythonScaffoldExecutor struct {
	PythonBin         string
	Timeout           time.Duration
	cfg               *config.Config
	pluginRepository  repository.PluginRepository
	projectRepository repository.ProjectRepository
}

type ScaffoldExecutionResult struct {
	RepoURL     string
	ProjectID   uuid.UUID
	ServiceName string
}

type scaffoldPluginPayload struct {
	Environment       string `json:"environment"`
	ServiceName       string `json:"service_name"`
	Port              int    `json:"port"`
	Database          string `json:"database"`
	ImageTag          string `json:"image_tag"`
	ModulePath        string `json:"module_path"`
	CIRegistryHost    string `json:"ci_registry_host"`
	CIServerURL       string `json:"ci_server_url"`
	CDProjectName     string `json:"cd_project_name"`
	CDRepoURL         string `json:"cd_repo_url"`
	CDTargetRevision  string `json:"cd_target_revision"`
	CDNamespace       string `json:"cd_namespace"`
	CDImageRepository string `json:"cd_image_repository"`
}

type scaffoldPluginInput struct {
	Action        string                `json:"action"`
	CorrelationID string                `json:"correlation_id"`
	Payload       scaffoldPluginPayload `json:"payload"`
}

type scaffoldPluginOutput struct {
	Status string `json:"status"`
	Output struct {
		RepoURL string `json:"repo_url"`
		Path    string `json:"path"`
	} `json:"output"`
}

var _ core.Executor[ScaffoldJob, ScaffoldExecutionResult] = (*ScaffoldExecutorAdapter)(nil)

func NewPythonScaffoldExecutor(
	cfg *config.Config,
	pluginRepository repository.PluginRepository,
	projectRepository repository.ProjectRepository,
) *PythonScaffoldExecutor {
	return &PythonScaffoldExecutor{
		PythonBin:         "python3",
		cfg:               cfg,
		pluginRepository:  pluginRepository,
		projectRepository: projectRepository,
		Timeout:           5 * time.Minute,
	}
}

func (e *PythonScaffoldExecutor) Execute(ctx context.Context, job *ScaffoldJob) (ScaffoldExecutionResult, error) {
	if job == nil {
		return ScaffoldExecutionResult{}, errors.New("job is nil")
	}

	if e.pluginRepository == nil {
		return ScaffoldExecutionResult{}, errors.New("plugin repository is required")
	}

	if e.projectRepository == nil {
		return ScaffoldExecutionResult{}, errors.New("project repository is required")
	}

	if e.cfg == nil {
		return ScaffoldExecutionResult{}, errors.New("config is required")
	}

	plugin, err := e.pluginRepository.FindOne(ctx, job.PluginID)

	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
			return ScaffoldExecutionResult{}, misc.WrapError(err, errs.NewInternalServerError("failed to find plugin by ID", nil))
		}
		return ScaffoldExecutionResult{}, err // Return the NotFoundError directly
	}

	project, err := e.projectRepository.FindOne(ctx, job.ProjectID)
	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) {
			return ScaffoldExecutionResult{}, misc.WrapError(err, errs.NewInternalServerError("failed to find project by ID", nil))
		}
		return ScaffoldExecutionResult{}, err
	}
	if project == nil {
		return ScaffoldExecutionResult{}, errors.New("project is required")
	}

	scaffoldRepoURL, _ := buildScaffoldRepoURL(
		strings.TrimSpace(e.cfg.ScmConfig.ExternalURL),
		project.OwnerTeam,
		job.Variables.ServiceName,
		project.ScmProvider,
	)

	CDRepoURL, _ := buildCDRepoURL(
		strings.TrimSpace(e.cfg.ArgoCD.RepoURL),
		project.OwnerTeam,
		job.Variables.ServiceName,
	)

	scriptPath := strings.TrimSpace(plugin.Entrypoint)

	if scriptPath == "" {
		return ScaffoldExecutionResult{}, errors.New("script path is required")
	}

	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

	payload := scaffoldPluginPayload{
		Environment:       job.Environment.String(),
		ServiceName:       job.Variables.ServiceName,
		Port:              job.Variables.Port,
		Database:          job.Variables.Database,
		ImageTag:          DEFAULT_IMAGE_TAG,
		ModulePath:        job.Variables.ModulePath,
		CIRegistryHost:    strings.TrimSpace(e.cfg.CI.ImageRegistryHost),
		CIServerURL:       strings.TrimSpace(e.cfg.CI.ServerURL),
		CDProjectName:     strings.TrimSpace(e.cfg.ArgoCD.AppProject),
		CDRepoURL:         CDRepoURL,
		CDTargetRevision:  strings.TrimSpace(e.cfg.ArgoCD.TargetRevision),
		CDNamespace:       strings.TrimSpace(e.cfg.ArgoCD.AppNamespace),
		CDImageRepository: strings.TrimSpace(e.cfg.ArgoCD.RepositoryRegistryHost) + "/" + job.Variables.ServiceName,
	}

	// TODO: use enum instead
	in := scaffoldPluginInput{
		Action:        "scaffold",
		CorrelationID: job.ID.String(),
		Payload:       payload,
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

	var out scaffoldPluginOutput

	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return ScaffoldExecutionResult{}, fmt.Errorf("invalid scaffold json output: %w", err)
	}

	if strings.ToLower(out.Status) != "ok" {
		return ScaffoldExecutionResult{}, fmt.Errorf("plugin returned non-ok status")
	}

	if scaffoldRepoURL == "" {
		return ScaffoldExecutionResult{}, errors.New("plugin output missing repo_url/path")
	}

	return ScaffoldExecutionResult{RepoURL: scaffoldRepoURL, ProjectID: job.ProjectID, ServiceName: job.Variables.ServiceName}, nil
}

func buildScaffoldRepoURL(baseURL string, owner string, serviceName string, scmProvider string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	owner = strings.TrimSpace(owner)
	serviceName = strings.TrimSpace(serviceName)

	if baseURL == "" {
		return "", errors.New("scm external url is required")
	}
	if owner == "" {
		return "", errors.New("project owner team is required")
	}
	if serviceName == "" {
		return "", errors.New("service name is required")
	}

	switch strings.ToLower(strings.TrimSpace(scmProvider)) {
	case "", "gitea", "github", "gitlab":
	default:
		return "", fmt.Errorf("unsupported scm provider %q", scmProvider)
	}

	parsed, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "", err
	}

	parsed.Path = path.Join(parsed.Path, owner, serviceName+".git")

	return parsed.String(), nil
}

func buildCDRepoURL(baseURL string, owner string, serviceName string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	owner = strings.TrimSpace(owner)
	serviceName = strings.TrimSpace(serviceName)

	if baseURL == "" {
		return "", errors.New("cd base url is required")
	}
	if owner == "" {
		return "", errors.New("project owner team is required")
	}
	if serviceName == "" {
		return "", errors.New("service name is required")
	}

	parsed, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "", err
	}

	parsed.Path = path.Join(parsed.Path, owner, serviceName+".git")

	return parsed.String(), nil
}
