package deployment

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"
)

type QueueSourceAdapter struct {
	deploymentRepository repository.DeploymentRepository
}

var _ core.QueueSourceAdapter[DeploymentJob] = (*QueueSourceAdapter)(nil)

func NewQueueSourceAdapter(deploymentRepository repository.DeploymentRepository) *QueueSourceAdapter {
	return &QueueSourceAdapter{deploymentRepository: deploymentRepository}
}

func (a *QueueSourceAdapter) Dequeue(ctx context.Context) (*DeploymentJob, error) {
	deployment, err := a.deploymentRepository.FindOnePending(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeue deployment: %w", err)
	}
	if deployment == nil {
		return nil, nil
	}

	return &DeploymentJob{Deployment: *deployment}, nil
}
