package deployment

import (
	"bytes"
	"context"
	"devhub-backend/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/domain/repository"
	"devhub-backend/internal/util/misc"
)

type CommandExecutor struct {
	Server            string
	AuthToken         string
	Insecure          bool
	Timeout           time.Duration
	RepoBaseURL       string
	AutoCreateApp     bool
	AppProject        string
	AppNamespace      string
	AppDestServer     string
	projectRepository repository.ProjectRepository
}

type ExecutionResult struct {
	ExternalRef string
	CommitSHA   string
	FinishedAt  time.Time
}

const defaultArgoCDPath = "k8s"

func NewCommandExecutor(cfg config.ArgoCDConfig, projectRepository repository.ProjectRepository) *CommandExecutor {
	return &CommandExecutor{
		Server:            strings.TrimSpace(cfg.Server),
		AuthToken:         strings.TrimSpace(cfg.AuthToken),
		Insecure:          cfg.Insecure,
		Timeout:           cfg.Timeout,
		RepoBaseURL:       strings.TrimSpace(cfg.RepoBaseURL),
		AutoCreateApp:     cfg.AutoCreateApp,
		AppProject:        strings.TrimSpace(cfg.AppProject),
		AppNamespace:      strings.TrimSpace(cfg.AppNamespace),
		AppDestServer:     strings.TrimSpace(cfg.AppDestServer),
		projectRepository: projectRepository,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, job *DeploymentJob) (ExecutionResult, error) {
	if job == nil {
		return ExecutionResult{}, errors.New("job is nil")
	}

	if e.Server == "" {
		return ExecutionResult{}, fmt.Errorf("%s is not configured", config.ArgoCDServerKey)
	}

	if e.AuthToken == "" {
		return ExecutionResult{}, fmt.Errorf("%s is not configured", config.ArgoCDAuthTokenKey)
	}

	if e.projectRepository == nil {
		return ExecutionResult{}, errors.New("project repository is required")
	}

	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

	project, err := e.projectRepository.FindOne(ctx, job.ProjectID)
	if err != nil {
		if !errors.As(err, &errs.NotFoundError{}) {
			return ExecutionResult{}, misc.WrapError(err, errs.NewInternalServerError("failed to find project by ID", nil))
		}
		return ExecutionResult{}, err
	}

	if project == nil {
		return ExecutionResult{}, errors.New("project is required")
	}

	if strings.TrimSpace(project.RepoURL) == "" {
		return ExecutionResult{}, errors.New("project repo_url is required")
	}

	repoURL := rewriteRepoURLForArgoCD(project.RepoURL, e.RepoBaseURL)
	revision := strings.TrimSpace(job.Version)
	if revision == "" {
		return ExecutionResult{}, errors.New("deployment revision is required")
	}
	args := []string{
		"--server",
		e.Server,
		"--auth-token",
		e.AuthToken,
	}

	if e.Insecure {
		args = append(args, "--insecure")
	}

	if err := e.ensureApplicationExists(ctx, args, job.Service, project.RepoURL, repoURL, revision); err != nil {
		return ExecutionResult{}, err
	}

	setArgs := append([]string{}, args...)
	setStdout, setStderr, err := runSetCommand(ctx, setArgs, job.Service, repoURL, revision)
	if err != nil {
		return ExecutionResult{}, formatArgoCDError("set", job.Service, project.RepoURL, repoURL, revision, err, setStdout, setStderr)
	}

	syncArgs := append([]string{}, args...)
	stdout, stderr, err := runSyncCommand(ctx, syncArgs, job.Service, revision)
	if err != nil {
		return ExecutionResult{}, formatArgoCDError("sync", job.Service, project.RepoURL, repoURL, revision, err, stdout, stderr)
	}

	result := ExecutionResult{
		FinishedAt: time.Now().UTC(),
	}

	if err := enrichExecutionResultFromApp(ctx, args, job.Service, &result); err != nil {
		return ExecutionResult{}, err
	}

	return result, nil
}

func (e *CommandExecutor) ensureApplicationExists(
	ctx context.Context,
	baseArgs []string,
	appName string,
	projectRepoURL string,
	argoRepoURL string,
	revision string,
) error {
	stdout, stderr, err := runGetCommand(ctx, baseArgs, appName)
	if err == nil {
		return nil
	}

	if !isArgoCDAppNotFound(stderr) {
		return fmt.Errorf(
			"check argocd application existence failed: app=%q stdout=%s stderr=%s: %w",
			appName,
			strings.TrimSpace(stdout),
			strings.TrimSpace(stderr),
			err,
		)
	}

	if !e.AutoCreateApp {
		return fmt.Errorf(
			"argocd application %q does not exist. Create it first or enable %s to allow automatic creation. project_repo=%q argocd_repo=%q revision=%q stderr=%s",
			appName,
			config.ArgoCDAutoCreateAppKey,
			projectRepoURL,
			argoRepoURL,
			revision,
			strings.TrimSpace(stderr),
		)
	}

	createStdout, createStderr, createErr := runCreateCommand(
		ctx,
		baseArgs,
		appName,
		argoRepoURL,
		revision,
		e.AppProject,
		e.AppNamespace,
		e.AppDestServer,
	)
	if createErr != nil {
		return formatArgoCDError("create", appName, projectRepoURL, argoRepoURL, revision, createErr, createStdout, createStderr)
	}

	return nil
}

func enrichExecutionResultFromApp(ctx context.Context, baseArgs []string, appName string, result *ExecutionResult) error {
	getArgs := append([]string{}, baseArgs...)
	getArgs = append(getArgs,
		"app",
		"get",
		appName,
		"-o",
		"json",
	)

	stdout, stderr, err := runArgoCDCommand(ctx, getArgs)
	if err != nil {
		return fmt.Errorf(
			"fetch argocd application state failed: %w; app=%q stdout=%s stderr=%s",
			err,
			appName,
			strings.TrimSpace(stdout),
			strings.TrimSpace(stderr),
		)
	}

	var parsed struct {
		Metadata struct {
			Name string `json:"name"`
			UID  string `json:"uid"`
		} `json:"metadata"`
		Status struct {
			Sync struct {
				Revision string `json:"revision"`
			} `json:"sync"`
			OperationState struct {
				SyncResult struct {
					Revision string `json:"revision"`
				} `json:"syncResult"`
			} `json:"operationState"`
		} `json:"status"`
	}

	if err := json.Unmarshal([]byte(stdout), &parsed); err != nil {
		return fmt.Errorf("decode argocd application json failed: %w", err)
	}

	result.ExternalRef = strings.TrimSpace(parsed.Metadata.UID)
	if result.ExternalRef == "" {
		result.ExternalRef = strings.TrimSpace(parsed.Metadata.Name)
	}

	result.CommitSHA = strings.TrimSpace(parsed.Status.Sync.Revision)
	if result.CommitSHA == "" {
		result.CommitSHA = strings.TrimSpace(parsed.Status.OperationState.SyncResult.Revision)
	}

	return nil
}

func runSetCommand(ctx context.Context, baseArgs []string, appName string, repoURL string, revision string) (string, string, error) {
	args := append([]string{}, baseArgs...)
	args = append(args,
		"app",
		"set",
		appName,
		"--repo",
		repoURL,
		"--revision",
		revision,
		"--path",
		defaultArgoCDPath,
	)
	return runArgoCDCommand(ctx, args)
}

func runGetCommand(ctx context.Context, baseArgs []string, appName string) (string, string, error) {
	args := append([]string{}, baseArgs...)
	args = append(args,
		"app",
		"get",
		appName,
	)
	return runArgoCDCommand(ctx, args)
}

func runCreateCommand(
	ctx context.Context,
	baseArgs []string,
	appName string,
	repoURL string,
	revision string,
	project string,
	namespace string,
	destServer string,
) (string, string, error) {
	args := append([]string{}, baseArgs...)
	args = append(args,
		"app",
		"create",
		appName,
		"--repo",
		repoURL,
		"--revision",
		revision,
		"--path",
		defaultArgoCDPath,
		"--project",
		project,
		"--dest-namespace",
		namespace,
		"--dest-server",
		destServer,
	)
	return runArgoCDCommand(ctx, args)
}

func runSyncCommand(ctx context.Context, baseArgs []string, appName string, revision string) (string, string, error) {
	args := append([]string{}, baseArgs...)
	args = append(args,
		"app",
		"sync",
		appName,
		"--revision",
		revision,
	)
	return runArgoCDCommand(ctx, args)
}

func runArgoCDCommand(ctx context.Context, args []string) (string, string, error) {
	cmd := exec.CommandContext(ctx, "argocd", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func rewriteRepoURLForArgoCD(repoURL string, repoBaseURL string) string {
	repoURL = strings.TrimSpace(repoURL)
	repoBaseURL = strings.TrimSpace(repoBaseURL)

	if repoURL == "" || repoBaseURL == "" {
		return repoURL
	}

	repoParsed, err := url.Parse(repoURL)
	if err != nil || repoParsed.Path == "" {
		return repoURL
	}

	baseParsed, err := url.Parse(repoBaseURL)
	if err != nil || baseParsed.Scheme == "" || baseParsed.Host == "" {
		return repoURL
	}

	rewritten := url.URL{
		Scheme:   baseParsed.Scheme,
		Host:     baseParsed.Host,
		Path:     repoParsed.Path,
		RawQuery: repoParsed.RawQuery,
		Fragment: repoParsed.Fragment,
	}
	return rewritten.String()
}

func formatArgoCDError(
	operation string,
	appName string,
	projectRepoURL string,
	argoRepoURL string,
	revision string,
	err error,
	stdout string,
	stderr string,
) error {
	stdout = strings.TrimSpace(stdout)
	stderr = strings.TrimSpace(stderr)
	lowerStderr := strings.ToLower(stderr)

	contextLine := fmt.Sprintf(
		"app=%q project_repo=%q argocd_repo=%q revision=%q",
		appName,
		projectRepoURL,
		argoRepoURL,
		revision,
	)

	if strings.Contains(lowerStderr, "permissiondenied") || strings.Contains(lowerStderr, "permission denied") {
		if operation == "set" {
			return fmt.Errorf(
				"argocd app set was denied. The configured token likely lacks Argo CD RBAC permission to update applications. Typical required permissions: applications,get and applications,update. %s stderr=%s",
				contextLine,
				stderr,
			)
		}
		if operation == "create" {
			return fmt.Errorf(
				"argocd app create was denied. The configured token likely lacks Argo CD RBAC permission to create applications. Typical required permissions: applications,get and applications,create. %s stderr=%s",
				contextLine,
				stderr,
			)
		}
		return fmt.Errorf(
			"argocd app sync was denied. The configured token likely lacks Argo CD RBAC permission to sync applications. Typical required permissions: applications,get and applications,sync. %s stderr=%s",
			contextLine,
			stderr,
		)
	}

	if strings.Contains(lowerStderr, "repository not accessible") || strings.Contains(lowerStderr, "unable to ls-remote") {
		switch {
		case strings.Contains(lowerStderr, "no such host"), strings.Contains(lowerStderr, "lookup "):
			return fmt.Errorf(
				"argocd could not resolve the repository host while validating the application source. Check cluster DNS and ARGOCD_REPO_BASE_URL. %s stderr=%s",
				contextLine,
				stderr,
			)
		case strings.Contains(lowerStderr, "connection refused"), strings.Contains(lowerStderr, "connect: connection refused"):
			return fmt.Errorf(
				"argocd reached the repository address but nothing accepted the connection. The rewritten repo base is likely wrong for the cluster network. Check that the Gitea HTTP endpoint is exposed and that ARGOCD_REPO_BASE_URL points to a reachable host:port. %s stderr=%s",
				contextLine,
				stderr,
			)
		case strings.Contains(lowerStderr, "authentication required"), strings.Contains(lowerStderr, "authorization failed"), strings.Contains(lowerStderr, "access denied"):
			return fmt.Errorf(
				"argocd could reach the repository but could not authenticate. Check repository credentials configured in Argo CD for this repo. %s stderr=%s",
				contextLine,
				stderr,
			)
		default:
			return fmt.Errorf(
				"argocd reported that the repository is not accessible while validating the application source. %s stderr=%s",
				contextLine,
				stderr,
			)
		}
	}

	if strings.Contains(lowerStderr, "unable to resolve") && strings.Contains(lowerStderr, "to a commit sha") {
		return fmt.Errorf(
			"argocd could access the repository but the requested revision does not exist as an exact Git ref. Provide a real branch, tag, or commit SHA such as \"main\" or \"v1.0.0\". %s stderr=%s",
			contextLine,
			stderr,
		)
	}

	return fmt.Errorf(
		"execute deployment %s command failed: %w; %s stdout=%s stderr=%s",
		operation,
		err,
		contextLine,
		stdout,
		stderr,
	)
}

func isArgoCDAppNotFound(stderr string) bool {
	lower := strings.ToLower(strings.TrimSpace(stderr))
	return strings.Contains(lower, "not found") || strings.Contains(lower, "does not exist")
}
