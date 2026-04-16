package release

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
)

type QueueSourceAdapter struct {
	releaseRepository repository.ReleaseRepository
}

var _ core.QueueSourceAdapter[ReleaseJob] = (*QueueSourceAdapter)(nil)

func NewQueueSourceAdapter(releaseRepository repository.ReleaseRepository) *QueueSourceAdapter {
	return &QueueSourceAdapter{releaseRepository: releaseRepository}
}

func (a *QueueSourceAdapter) Dequeue(ctx context.Context) (*ReleaseJob, error) {
	release, err := a.releaseRepository.FindOnePending(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeue release: %w", err)
	}
	if release == nil {
		return nil, nil
	}

	return &ReleaseJob{Release: *release}, nil
}
