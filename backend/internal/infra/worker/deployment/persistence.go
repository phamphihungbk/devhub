package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
	"devhub-backend/internal/util/misc"

	"github.com/google/uuid"
)

type StatePersistence struct {
	jobRepository        repository.JobRepository
	deploymentRepository repository.DeploymentRepository
}

var _ core.StatePersistence[DeploymentExecutionResult] = (*StatePersistence)(nil)

func NewStatePersistence(jobRepository repository.JobRepository, deploymentRepository repository.DeploymentRepository) *StatePersistence {
	return &StatePersistence{jobRepository: jobRepository, deploymentRepository: deploymentRepository}
}

func (p *StatePersistence) MarkRunning(ctx context.Context, id uuid.UUID) error {
	job, err := p.jobRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find deployment job before marking running: %w", err)
	}
	if job == nil {
		return fmt.Errorf("deployment job %s not found", id)
	}
	if job.Status != entity.JobStatusQueued {
		return fmt.Errorf("deployment job %s is not queued", id)
	}

	deployment, err := p.deploymentRepository.FindOne(ctx, job.ResourceID)
	if err != nil {
		return fmt.Errorf("find deployment before marking running: %w", err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s not found", job.ResourceID)
	}
	// if deployment.Status != entity.DeploymentStatusPending {
	// 	return fmt.Errorf("deployment %s is not pending", job.ResourceID)
	// }

	now := time.Now().UTC()
	jobStatus := entity.JobStatusRunning
	attempts := job.Attempts + 1
	if _, err := p.jobRepository.UpdateOne(ctx, repository.UpdateJobInput{
		ID:        id,
		Status:    &jobStatus,
		Attempts:  &attempts,
		StartedAt: &now,
	}); err != nil {
		return fmt.Errorf("mark deployment job running: %w", err)
	}

	status := entity.DeploymentStatusRunning
	if _, err := p.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:     job.ResourceID,
		Status: &status,
	}); err != nil {
		return fmt.Errorf("mark deployment running: %w", err)
	}

	return nil
}

func (p *StatePersistence) MarkCompleted(ctx context.Context, id uuid.UUID, result DeploymentExecutionResult) error {
	job, err := p.jobRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find deployment job before marking completed: %w", err)
	}
	if job == nil {
		return fmt.Errorf("deployment job %s not found", id)
	}
	if job.Status != entity.JobStatusRunning {
		return fmt.Errorf("deployment job %s is not running", id)
	}

	deployment, err := p.deploymentRepository.FindOne(ctx, job.ResourceID)
	if err != nil {
		return fmt.Errorf("find deployment before marking completed: %w", err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s not found", job.ResourceID)
	}
	if deployment.Status != entity.DeploymentStatusRunning {
		return fmt.Errorf("deployment %s is not running", job.ResourceID)
	}

	finishedAt := time.Now()
	status := entity.DeploymentStatusCompleted

	if _, err := p.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:          job.ResourceID,
		Status:      misc.ToPointer(status),
		ExternalRef: misc.ToPointer(result.ExternalRef),
		CommitSHA:   misc.ToPointer(result.CommitSHA),
		FinishedAt:  misc.ToPointer(finishedAt),
	}); err != nil {
		return fmt.Errorf("mark deployment completed: %w", err)
	}

	resultBytes, err := json.Marshal(map[string]string{
		"external_ref": result.ExternalRef,
		"commit_sha":   result.CommitSHA,
	})
	if err != nil {
		return fmt.Errorf("marshal deployment job result: %w", err)
	}
	jobStatus := entity.JobStatusCompleted
	resultJSON := string(resultBytes)
	if _, err := p.jobRepository.UpdateOne(ctx, repository.UpdateJobInput{
		ID:         id,
		Status:     misc.ToPointer(jobStatus),
		Result:     misc.ToPointer(resultJSON),
		FinishedAt: misc.ToPointer(finishedAt),
	}); err != nil {
		return fmt.Errorf("mark deployment job completed: %w", err)
	}

	return nil
}

func (p *StatePersistence) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	job, err := p.jobRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find deployment job before marking failed: %w", err)
	}
	if job == nil {
		return fmt.Errorf("deployment job %s not found", id)
	}
	if job.Status != entity.JobStatusRunning {
		return fmt.Errorf("deployment job %s is not running", id)
	}

	deployment, err := p.deploymentRepository.FindOne(ctx, job.ResourceID)
	if err != nil {
		return fmt.Errorf("find deployment before marking failed: %w", err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s not found", job.ResourceID)
	}
	if deployment.Status != entity.DeploymentStatusRunning {
		return fmt.Errorf("deployment %s is not running", job.ResourceID)
	}

	status := entity.DeploymentStatusFailed
	finishedAt := time.Now().UTC()

	if _, err := p.deploymentRepository.UpdateOne(ctx, repository.UpdateDeploymentInput{
		ID:         job.ResourceID,
		Status:     misc.ToPointer(status),
		FinishedAt: misc.ToPointer(finishedAt),
	}); err != nil {
		return fmt.Errorf("mark deployment failed: %w", err)
	}

	jobStatus := entity.JobStatusFailed
	if _, err := p.jobRepository.UpdateOne(ctx, repository.UpdateJobInput{
		ID:         id,
		Status:     misc.ToPointer(jobStatus),
		Error:      misc.ToPointer(reason),
		FinishedAt: misc.ToPointer(finishedAt),
	}); err != nil {
		return fmt.Errorf("mark deployment job failed: %w", err)
	}

	return nil
}
