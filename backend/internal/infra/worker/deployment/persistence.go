package deployment

import (
	"context"
	"fmt"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

type StatePersistence struct {
	deploymentRepository repository.DeploymentRepository
}

var _ core.StatePersistence[DeploymentExecutionResult] = (*StatePersistence)(nil)

func NewStatePersistence(deploymentRepository repository.DeploymentRepository) *StatePersistence {
	return &StatePersistence{deploymentRepository: deploymentRepository}
}

func (p *StatePersistence) MarkRunning(ctx context.Context, id uuid.UUID) error {
	deployment, err := p.deploymentRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find deployment before marking running: %w", err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s not found", id)
	}
	if deployment.Status != entity.DeploymentStatusPending {
		return fmt.Errorf("deployment %s is not pending", id)
	}

	status := entity.DeploymentStatusRunning
	if _, err := p.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:     id,
		Status: &status,
	}); err != nil {
		return fmt.Errorf("mark deployment running: %w", err)
	}

	return nil
}

func (p *StatePersistence) MarkCompleted(ctx context.Context, id uuid.UUID, result DeploymentExecutionResult) error {
	deployment, err := p.deploymentRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find deployment before marking completed: %w", err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s not found", id)
	}
	if deployment.Status != entity.DeploymentStatusRunning {
		return fmt.Errorf("deployment %s is not running", id)
	}

	status := entity.DeploymentStatusCompleted
	finishedAt := result.FinishedAt
	if finishedAt.IsZero() {
		finishedAt = time.Now().UTC()
	}

	if _, err := p.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:           id,
		Status:       &status,
		ExternalRef:  optionalString(result.ExternalRef),
		CommitSHA:    optionalString(result.CommitSHA),
		RunnerOutput: optionalString(result.RunnerOutput),
		RunnerError:  optionalString(result.RunnerError),
		FinishedAt:   &finishedAt,
	}); err != nil {
		return fmt.Errorf("mark deployment completed: %w", err)
	}

	return nil
}

func (p *StatePersistence) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	deployment, err := p.deploymentRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find deployment before marking failed: %w", err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s not found", id)
	}
	if deployment.Status != entity.DeploymentStatusRunning {
		return fmt.Errorf("deployment %s is not running", id)
	}

	status := entity.DeploymentStatusFailed
	finishedAt := time.Now().UTC()

	if _, err := p.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:          id,
		Status:      &status,
		RunnerError: optionalString(reason),
		FinishedAt:  &finishedAt,
	}); err != nil {
		return fmt.Errorf("mark deployment failed: %w", err)
	}

	return nil
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
