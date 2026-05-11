package scaffold

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
)

type ScaffoldQueueSourceAdapter struct {
	jobRepository repository.JobRepository
}

var _ core.QueueSourceAdapter[ScaffoldJob] = (*ScaffoldQueueSourceAdapter)(nil)

func NewScaffoldQueueSourceAdapter(jobRepository repository.JobRepository) *ScaffoldQueueSourceAdapter {
	return &ScaffoldQueueSourceAdapter{jobRepository: jobRepository}
}

func (a *ScaffoldQueueSourceAdapter) Dequeue(ctx context.Context) (*ScaffoldJob, error) {
	job, err := a.jobRepository.FindOnePending(ctx, entity.JobTypeScaffold)
	if err != nil {
		return nil, fmt.Errorf("dequeue scaffold job: %w", err)
	}
	if job == nil {
		return nil, nil
	}

	return &ScaffoldJob{
		Job: *job,
	}, nil
}
