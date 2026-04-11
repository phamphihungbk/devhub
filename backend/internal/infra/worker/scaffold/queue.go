package scaffold

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
)

type ScaffoldQueueSourceAdapter struct {
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

var _ core.QueueSourceAdapter[ScaffoldJob] = (*ScaffoldQueueSourceAdapter)(nil)

func NewScaffoldQueueSourceAdapter(scaffoldRequestRepository repository.ScaffoldRequestRepository) *ScaffoldQueueSourceAdapter {
	return &ScaffoldQueueSourceAdapter{scaffoldRequestRepository: scaffoldRequestRepository}
}

func (a *ScaffoldQueueSourceAdapter) Dequeue(ctx context.Context) (*ScaffoldJob, error) {
	scaffoldRequest, err := a.scaffoldRequestRepository.FindOnePending(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeue scaffold request: %w", err)
	}
	if scaffoldRequest == nil {
		return nil, nil
	}

	return &ScaffoldJob{
		ScaffoldRequest: *scaffoldRequest,
	}, nil
}
