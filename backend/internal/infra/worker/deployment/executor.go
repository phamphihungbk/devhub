package deployment

import (
	"bytes"
	"context"
	"devhub-backend/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	AutoBuildImage    bool
	ImageBuilder      string
	MinikubeProfile   string
	GiteaURL          string
	GiteaExternalURL  string
	projectRepository repository.ProjectRepository
}

type ExecutionResult struct {
	ExternalRef string
	CommitSHA   string
	FinishedAt  time.Time
}

const defaultArgoCDPath = "charts/app"

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
		AutoBuildImage:    cfg.AutoBuildImage,
		ImageBuilder:      strings.TrimSpace(cfg.ImageBuilder),
		MinikubeProfile:   strings.TrimSpace(cfg.MinikubeProfile),
		projectRepository: projectRepository,
	}
}

func NewCommandExecutorWithGitProvider(cfg config.ArgoCDConfig, giteaCfg config.GiteaConfig, projectRepository repository.ProjectRepository) *CommandExecutor {
	executor := NewCommandExecutor(cfg, projectRepository)
	executor.GiteaURL = strings.TrimSpace(giteaCfg.URL)
	executor.GiteaExternalURL = strings.TrimSpace(giteaCfg.ExternalURL)
	return executor
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

	if err := e.ensureImageBuilt(ctx, project.RepoURL, revision); err != nil {
		return ExecutionResult{}, err
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

var (
	imageLinePattern        = regexp.MustCompile(`(?m)^\s*image:\s*["']?([^"'\s#]+)["']?`)
	imageRepositoryPattern  = regexp.MustCompile(`(?m)^\s*repository:\s*["']?([^"'\s#]+)["']?`)
	imageTagPattern         = regexp.MustCompile(`(?m)^\s*tag:\s*["']?([^"'\s#]+)["']?`)
)

func (e *CommandExecutor) ensureImageBuilt(ctx context.Context, projectRepoURL string, revision string) error {
	if !e.AutoBuildImage {
		return nil
	}

	cloneURL := rewriteRepoURLForWorker(projectRepoURL, e.GiteaURL, e.GiteaExternalURL)
	repoDir, cleanup, err := checkoutRepository(ctx, cloneURL, revision)
	if err != nil {
		return err
	}
	defer cleanup()

	image, err := readImageForBuild(repoDir)
	if err != nil {
		return err
	}

	if err := buildContainerImage(ctx, e.ImageBuilder, e.MinikubeProfile, repoDir, image); err != nil {
		return err
	}

	return nil
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

func checkoutRepository(ctx context.Context, repoURL string, revision string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "devhub-deploy-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temporary repository directory failed: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	cloneCmd := exec.CommandContext(ctx, "git", "clone", repoURL, tmpDir)
	var cloneStdout bytes.Buffer
	var cloneStderr bytes.Buffer
	cloneCmd.Stdout = &cloneStdout
	cloneCmd.Stderr = &cloneStderr

	if err := cloneCmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf(
			"clone deployment repository failed: repo=%q revision=%q stdout=%s stderr=%s: %w",
			repoURL,
			revision,
			strings.TrimSpace(cloneStdout.String()),
			strings.TrimSpace(cloneStderr.String()),
			err,
		)
	}

	checkoutCmd := exec.CommandContext(ctx, "git", "-C", tmpDir, "checkout", revision)
	var checkoutStdout bytes.Buffer
	var checkoutStderr bytes.Buffer
	checkoutCmd.Stdout = &checkoutStdout
	checkoutCmd.Stderr = &checkoutStderr

	if err := checkoutCmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf(
			"checkout deployment revision failed: repo=%q revision=%q stdout=%s stderr=%s: %w",
			repoURL,
			revision,
			strings.TrimSpace(checkoutStdout.String()),
			strings.TrimSpace(checkoutStderr.String()),
			err,
		)
	}

	return tmpDir, cleanup, nil
}

func readImageForBuild(repoDir string) (string, error) {
	valuesPath := filepath.Join(repoDir, defaultArgoCDPath, "values.yaml")
	if image, err := readImageFromHelmValues(valuesPath); err == nil {
		return image, nil
	}

	manifestPath := filepath.Join(repoDir, "k8s", "deployment.yaml")
	return readImageFromManifest(manifestPath)
}

func readImageFromHelmValues(path string) (string, error) {
	values, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read helm values for image build failed: path=%q: %w", path, err)
	}

	repositoryMatch := imageRepositoryPattern.FindSubmatch(values)
	tagMatch := imageTagPattern.FindSubmatch(values)
	if len(repositoryMatch) < 2 || len(tagMatch) < 2 {
		return "", fmt.Errorf("helm values %q do not contain image.repository and image.tag", path)
	}

	repository := strings.TrimSpace(string(repositoryMatch[1]))
	tag := strings.TrimSpace(string(tagMatch[1]))
	if repository == "" {
		return "", fmt.Errorf("helm values %q contain an empty image.repository", path)
	}

	if tag == "" {
		return repository, nil
	}

	return fmt.Sprintf("%s:%s", repository, tag), nil
}

func readImageFromManifest(path string) (string, error) {
	manifest, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read deployment manifest for image build failed: path=%q: %w", path, err)
	}

	match := imageLinePattern.FindSubmatch(manifest)
	if len(match) < 2 {
		return "", fmt.Errorf("deployment manifest %q does not contain a container image", path)
	}

	image := strings.TrimSpace(string(match[1]))
	if image == "" {
		return "", fmt.Errorf("deployment manifest %q contains an empty container image", path)
	}

	return image, nil
}

func buildContainerImage(ctx context.Context, builder string, minikubeProfile string, repoDir string, image string) error {
	builder = strings.ToLower(strings.TrimSpace(builder))
	if builder == "" {
		builder = "minikube"
	}

	var args []string
	command := ""

	switch builder {
	case "docker":
		command = "docker"
		args = []string{"build", "-t", image, "."}
	case "minikube":
		command = "minikube"
		args = []string{"image", "build"}
		if strings.TrimSpace(minikubeProfile) != "" {
			args = append(args, "-p", strings.TrimSpace(minikubeProfile))
		}
		args = append(args, "-t", image, ".")
	default:
		return fmt.Errorf("unsupported image builder %q; supported values are \"docker\" and \"minikube\"", builder)
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = repoDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf(
				"deployment image build requires %q, but it is not available in PATH. builder=%q repo_dir=%q image=%q",
				command,
				builder,
				repoDir,
				image,
			)
		}
		return fmt.Errorf(
			"build deployment image failed: builder=%q repo_dir=%q image=%q stdout=%s stderr=%s: %w",
			builder,
			repoDir,
			image,
			strings.TrimSpace(stdout.String()),
			strings.TrimSpace(stderr.String()),
			err,
		)
	}

	return nil
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

func rewriteRepoURLForWorker(repoURL string, internalBaseURL string, externalBaseURL string) string {
	repoURL = strings.TrimSpace(repoURL)
	internalBaseURL = strings.TrimSpace(internalBaseURL)
	externalBaseURL = strings.TrimSpace(externalBaseURL)

	if repoURL == "" || internalBaseURL == "" || externalBaseURL == "" {
		return repoURL
	}

	repoParsed, err := url.Parse(repoURL)
	if err != nil || repoParsed.Path == "" {
		return repoURL
	}

	externalParsed, err := url.Parse(externalBaseURL)
	if err != nil || externalParsed.Scheme == "" || externalParsed.Host == "" {
		return repoURL
	}

	if !strings.EqualFold(repoParsed.Host, externalParsed.Host) {
		return repoURL
	}

	internalParsed, err := url.Parse(internalBaseURL)
	if err != nil || internalParsed.Scheme == "" || internalParsed.Host == "" {
		return repoURL
	}

	rewritten := url.URL{
		Scheme:   internalParsed.Scheme,
		Host:     internalParsed.Host,
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
