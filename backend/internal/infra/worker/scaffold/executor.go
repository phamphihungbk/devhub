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
	ScaffoldRequestID string `json:"scaffold_request_id"`
	ProjectID         string `json:"project_id"`
	Environment       string `json:"environment"`
	ServiceName       string `json:"service_name"`
	Port              int    `json:"port"`
	Database          string `json:"database"`
	EnableLogging     bool   `json:"enable_logging"`
	RepoURL           string `json:"repo_url"`
	Namespace         string `json:"namespace"`
	TargetRevision    string `json:"target_revision"`
	ArgocdProject     string `json:"argocd_project"`
	RegistryURL       string `json:"registry_url"`
	ServerURL         string `json:"server_url"`
	ModulePath        string `json:"module_path"`
	Image             string `json:"image"`
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

	repoURL, err := buildScaffoldRepoURL(
		strings.TrimSpace(e.cfg.ArgoCD.RepoBaseURL),
		project.OwnerTeam,
		job.Variables.ServiceName,
		project.ScmProvider,
	)
	if err != nil {
		return ScaffoldExecutionResult{}, fmt.Errorf("build scaffold repo url: %w", err)
	}

	modulePath, err := inferModuleBaseFromRepoURL(repoURL)
	if err != nil {
		return ScaffoldExecutionResult{}, fmt.Errorf("infer module path from repo url: %w", err)
	}

	image := buildScaffoldImage(e.cfg.ArgoCD.ImageRegistryURL, job.Variables.ServiceName)

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
		ScaffoldRequestID: job.ID.String(),
		ProjectID:         job.ProjectID.String(),
		Environment:       job.Environment.String(),
		ServiceName:       job.Variables.ServiceName,
		Port:              job.Variables.Port,
		Database:          job.Variables.Database,
		EnableLogging:     job.Variables.EnableLogging,
		RepoURL:           repoURL,
		Namespace:         strings.TrimSpace(e.cfg.ArgoCD.AppNamespace),
		TargetRevision:    strings.TrimSpace(e.cfg.Gitops.Branch),
		ArgocdProject:     strings.TrimSpace(e.cfg.ArgoCD.AppProject),
		RegistryURL:       strings.TrimSpace(e.cfg.ArgoCD.ImageRegistryHost),
		ServerURL:         strings.TrimSpace(e.cfg.ArgoCD.RepoBaseURL),
		ModulePath:        modulePath,
		Image:             image,
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

	repoURL = strings.TrimSpace(out.Output.RepoURL)

	if repoURL == "" {
		repoURL = strings.TrimSpace(out.Output.Path)
	}

	if repoURL == "" {
		return ScaffoldExecutionResult{}, errors.New("plugin output missing repo_url/path")
	}

	return ScaffoldExecutionResult{RepoURL: repoURL, ProjectID: job.ProjectID, ServiceName: job.Variables.ServiceName}, nil
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

func inferModuleBaseFromRepoURL(repoURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(repoURL))
	if err != nil {
		return "", err
	}

	host := strings.TrimSpace(parsed.Host)
	repoPath := strings.Trim(parsed.Path, "/")
	if host == "" || repoPath == "" {
		return "", errors.New("repo url must include host and path")
	}

	segments := strings.Split(repoPath, "/")
	if len(segments) == 0 {
		return "", errors.New("repo url must include owner path")
	}

	ownerSegments := segments[:len(segments)-1]
	if len(ownerSegments) == 0 {
		return host, nil
	}

	return path.Join(append([]string{host}, ownerSegments...)...), nil
}

func buildScaffoldImage(registryURL string, serviceName string) string {
	registryURL = strings.TrimRight(strings.TrimSpace(registryURL), "/")
	serviceName = strings.TrimSpace(serviceName)

	if registryURL == "" {
		return serviceName + ":latest"
	}

	return registryURL + "/" + serviceName + ":latest"
}
