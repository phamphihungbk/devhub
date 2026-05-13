package scaffold

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

type ScaffoldStatePersistence struct {
	jobRepository             repository.JobRepository
	scaffoldRequestRepository repository.ScaffoldRequestRepository
	serviceRepository         repository.ServiceRepository
}

var _ core.StatePersistence[ScaffoldExecutionResult] = (*ScaffoldStatePersistence)(nil)

func NewScaffoldStatePersistence(
	jobRepository repository.JobRepository,
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
	serviceRepository repository.ServiceRepository,
) *ScaffoldStatePersistence {
	return &ScaffoldStatePersistence{jobRepository: jobRepository, scaffoldRequestRepository: scaffoldRequestRepository, serviceRepository: serviceRepository}
}

func (p *ScaffoldStatePersistence) MarkRunning(ctx context.Context, id uuid.UUID) error {
	job, err := p.jobRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find scaffold job before marking running: %w", err)
	}
	if job == nil {
		return fmt.Errorf("scaffold job %s not found", id)
	}
	if job.Status != entity.JobStatusQueued {
		return fmt.Errorf("scaffold job %s is not queued", id)
	}

	scaffoldRequest, err := p.scaffoldRequestRepository.FindOne(ctx, job.ResourceID)
	if err != nil {
		return fmt.Errorf("find scaffold request before marking running: %w", err)
	}
	if scaffoldRequest == nil {
		return fmt.Errorf("scaffold request %s not found", job.ResourceID)
	}
	// if scaffoldRequest.Status != entity.ScaffoldRequestApproved {
	// 	return fmt.Errorf("scaffold request %s is not approved", job.ResourceID)
	// }

	now := time.Now()
	jobStatus := entity.JobStatusRunning
	attempts := job.Attempts + 1
	if _, err := p.jobRepository.UpdateOne(ctx, repository.UpdateJobInput{
		ID:        id,
		Status:    &jobStatus,
		Attempts:  &attempts,
		StartedAt: &now,
	}); err != nil {
		return fmt.Errorf("mark scaffold job running: %w", err)
	}

	scaffoldStatus := entity.ScaffoldRequestRunning
	if _, err := p.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
		ID:     job.ResourceID,
		Status: &scaffoldStatus,
	}); err != nil {
		return fmt.Errorf("mark scaffold request running: %w", err)
	}

	return nil
}

func (p *ScaffoldStatePersistence) MarkCompleted(ctx context.Context, id uuid.UUID, result ScaffoldExecutionResult) error {
	job, err := p.jobRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find scaffold job before marking completed: %w", err)
	}
	if job == nil {
		return fmt.Errorf("scaffold job %s not found", id)
	}
	if job.Status != entity.JobStatusRunning {
		return fmt.Errorf("scaffold job %s is not running", id)
	}

	scaffoldRequest, err := p.scaffoldRequestRepository.FindOne(ctx, job.ResourceID)
	if err != nil {
		return fmt.Errorf("find scaffold request before marking completed: %w", err)
	}
	if scaffoldRequest == nil {
		return fmt.Errorf("scaffold request %s not found", job.ResourceID)
	}
	if scaffoldRequest.Status != entity.ScaffoldRequestRunning {
		return fmt.Errorf("scaffold request %s is not running", job.ResourceID)
	}

	scaffoldStatus := entity.ScaffoldRequestCompleted

	if _, err := p.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
		ID:            job.ResourceID,
		Status:        &scaffoldStatus,
		ResultRepoURL: &result.RepoURL,
	}); err != nil {
		return fmt.Errorf("mark scaffold request completed: %w", err)
	}

	resultBytes, err := json.Marshal(map[string]string{"repo_url": result.RepoURL})
	if err != nil {
		return fmt.Errorf("marshal scaffold job result: %w", err)
	}
	jobStatus := entity.JobStatusCompleted
	finishedAt := time.Now()
	resultJSON := string(resultBytes)
	if _, err := p.jobRepository.UpdateOne(ctx, repository.UpdateJobInput{
		ID:         id,
		Status:     &jobStatus,
		Result:     &resultJSON,
		FinishedAt: &finishedAt,
	}); err != nil {
		return fmt.Errorf("mark scaffold job completed: %w", err)
	}

	serviceName := job.Payload.VariableString("service_name")
	if serviceName == "" {
		return fmt.Errorf("job payload variable service_name is required")
	}

	repoURL := result.RepoURL
	if repoURL == "" {
		return fmt.Errorf("repo url is required")
	}

	if _, err := p.serviceRepository.CreateOne(ctx, &entity.Service{
		ProjectID: scaffoldRequest.ProjectID,
		Name:      serviceName,
		RepoURL:   repoURL,
		CreatedBy: scaffoldRequest.RequestedBy,
	}); err != nil {
		return err
	}

	return nil
}

func (p *ScaffoldStatePersistence) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	job, err := p.jobRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find scaffold job before marking failed: %w", err)
	}
	if job == nil {
		return fmt.Errorf("scaffold job %s not found", id)
	}
	if job.Status != entity.JobStatusRunning {
		return fmt.Errorf("scaffold job %s is not running", id)
	}

	scaffoldRequest, err := p.scaffoldRequestRepository.FindOne(ctx, job.ResourceID)
	if err != nil {
		return fmt.Errorf("find scaffold request before marking failed: %w", err)
	}
	if scaffoldRequest == nil {
		return fmt.Errorf("scaffold request %s not found", job.ResourceID)
	}
	if scaffoldRequest.Status != entity.ScaffoldRequestRunning {
		return fmt.Errorf("scaffold request %s is not running", job.ResourceID)
	}

	scaffoldStatus := entity.ScaffoldRequestFailed

	if _, err := p.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
		ID:     job.ResourceID,
		Status: &scaffoldStatus,
	}); err != nil {
		return fmt.Errorf("mark scaffold request failed: %w", err)
	}

	jobStatus := entity.JobStatusFailed
	finishedAt := time.Now()
	if _, err := p.jobRepository.UpdateOne(ctx, repository.UpdateJobInput{
		ID:         id,
		Status:     &jobStatus,
		Error:      &reason,
		FinishedAt: &finishedAt,
	}); err != nil {
		return fmt.Errorf("mark scaffold job failed: %w", err)
	}

	return nil
}
