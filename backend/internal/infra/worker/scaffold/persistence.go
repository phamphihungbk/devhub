package scaffold

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

type ScaffoldStatePersistence struct {
	scaffoldRequestRepository repository.ScaffoldRequestRepository
}

var _ core.StatePersistence[ScaffoldExecutionResult] = (*ScaffoldStatePersistence)(nil)

func NewScaffoldStatePersistence(scaffoldRequestRepository repository.ScaffoldRequestRepository) *ScaffoldStatePersistence {
	return &ScaffoldStatePersistence{scaffoldRequestRepository: scaffoldRequestRepository}
}

func (p *ScaffoldStatePersistence) MarkRunning(ctx context.Context, id uuid.UUID) error {
	scaffoldRequest, err := p.scaffoldRequestRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find scaffold request before marking running: %w", err)
	}
	if scaffoldRequest == nil {
		return fmt.Errorf("scaffold request %s not found", id)
	}
	if scaffoldRequest.Status != entity.ScaffoldRequestPending {
		return fmt.Errorf("scaffold request %s is not pending", id)
	}

	status := entity.ScaffoldRequestRunning
	if _, err := p.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
		ID:     id,
		Status: &status,
	}); err != nil {
		return fmt.Errorf("mark scaffold request running: %w", err)
	}

	return nil
}

func (p *ScaffoldStatePersistence) MarkCompleted(ctx context.Context, id uuid.UUID, result ScaffoldExecutionResult) error {
	scaffoldRequest, err := p.scaffoldRequestRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find scaffold request before marking completed: %w", err)
	}
	if scaffoldRequest == nil {
		return fmt.Errorf("scaffold request %s not found", id)
	}
	if scaffoldRequest.Status != entity.ScaffoldRequestRunning {
		return fmt.Errorf("scaffold request %s is not running", id)
	}

	status := entity.ScaffoldRequestCompleted

	if _, err := p.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
		ID:            id,
		Status:        &status,
		ResultRepoURL: &result.RepoURL,
	}); err != nil {
		return fmt.Errorf("mark scaffold request completed: %w", err)
	}
	return nil
}

func (p *ScaffoldStatePersistence) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	scaffoldRequest, err := p.scaffoldRequestRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find scaffold request before marking failed: %w", err)
	}
	if scaffoldRequest == nil {
		return fmt.Errorf("scaffold request %s not found", id)
	}
	if scaffoldRequest.Status != entity.ScaffoldRequestRunning {
		return fmt.Errorf("scaffold request %s is not running", id)
	}

	status := entity.ScaffoldRequestFailed

	if _, err := p.scaffoldRequestRepository.UpdateOne(ctx, repository.UpdateScaffoldRequestInput{
		ID:     id,
		Status: &status,
	}); err != nil {
		return fmt.Errorf("mark scaffold request failed: %w", err)
	}
	_ = reason
	return nil
}
