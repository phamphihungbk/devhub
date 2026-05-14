package deployment

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
)

type QueueSourceAdapter struct {
	jobRepository repository.JobRepository
}

var _ core.QueueSourceAdapter[DeploymentJob] = (*QueueSourceAdapter)(nil)

func NewQueueSourceAdapter(jobRepository repository.JobRepository) *QueueSourceAdapter {
	return &QueueSourceAdapter{jobRepository: jobRepository}
}

func (a *QueueSourceAdapter) Dequeue(ctx context.Context) (*DeploymentJob, error) {
	job, err := a.jobRepository.FindOnePending(ctx, entity.JobTypeDeployment)
	if err != nil {
		return nil, fmt.Errorf("dequeue deployment job: %w", err)
	}
	if job == nil {
		return nil, nil
	}

	return &DeploymentJob{Job: *job}, nil
}
