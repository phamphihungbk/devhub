package deployment

import (
	"bytes"
	"context"
	"devhub-backend/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CommandExecutor struct {
	Server    string
	AuthToken string
	Insecure  bool
	Timeout   time.Duration
}

type ExecutionResult struct {
	ExternalRef string
	CommitSHA   string
	FinishedAt  time.Time
}

func NewCommandExecutor(cfg config.ArgoCDConfig) *CommandExecutor {
	return &CommandExecutor{
		Server:    strings.TrimSpace(cfg.Server),
		AuthToken: strings.TrimSpace(cfg.AuthToken),
		Insecure:  cfg.Insecure,
		Timeout:   cfg.Timeout,
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

	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
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
	args = append(args,
		"app",
		"sync",
		job.Service,
		"--revision",
		job.Version,
	)

	stdout, stderr, err := runArgoCDCommand(ctx, args)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf(
			"execute deployment command: %w; stdout=%s; stderr=%s",
			err,
			strings.TrimSpace(stdout),
			strings.TrimSpace(stderr),
		)
	}

	result := ExecutionResult{
		FinishedAt: time.Now().UTC(),
	}

	if err := decodeExecutionResult([]byte(stdout), &result); err != nil {
		return ExecutionResult{}, err
	}

	return result, nil
}

func decodeExecutionResult(stdout []byte, result *ExecutionResult) error {
	trimmed := strings.TrimSpace(string(stdout))
	if trimmed == "" {
		return nil
	}

	var parsed struct {
		ExternalRef string `json:"external_ref"`
		CommitSHA   string `json:"commit_sha"`
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}

	result.ExternalRef = strings.TrimSpace(parsed.ExternalRef)
	result.CommitSHA = strings.TrimSpace(parsed.CommitSHA)
	return nil
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
