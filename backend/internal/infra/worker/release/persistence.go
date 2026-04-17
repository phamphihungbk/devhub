package release

import (
	"context"
	"fmt"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

type StatePersistence struct {
	releaseRepository repository.ReleaseRepository
}

var _ core.StatePersistence[ReleaseExecutionResult] = (*StatePersistence)(nil)

func NewStatePersistence(releaseRepository repository.ReleaseRepository) *StatePersistence {
	return &StatePersistence{releaseRepository: releaseRepository}
}

func (p *StatePersistence) MarkRunning(ctx context.Context, id uuid.UUID) error {
	release, err := p.releaseRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find release before marking running: %w", err)
	}
	if release == nil {
		return fmt.Errorf("release %s not found", id)
	}
	if release.Status != entity.ReleaseStatusPending {
		return fmt.Errorf("release %s is not pending", id)
	}

	status := entity.ReleaseStatusRunning
	if _, err := p.releaseRepository.UpdateOne(ctx, repository.UpdateReleaseInput{
		ID:     id,
		Status: &status,
	}); err != nil {
		return fmt.Errorf("mark release running: %w", err)
	}

	return nil
}

func (p *StatePersistence) MarkCompleted(ctx context.Context, id uuid.UUID, result ReleaseExecutionResult) error {
	release, err := p.releaseRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find release before marking completed: %w", err)
	}
	if release == nil {
		return fmt.Errorf("release %s not found", id)
	}
	if release.Status != entity.ReleaseStatusRunning {
		return fmt.Errorf("release %s is not running", id)
	}

	status := entity.ReleaseStatusCompleted
	if _, err := p.releaseRepository.UpdateOne(ctx, repository.UpdateReleaseInput{
		ID:          id,
		Status:      &status,
		ExternalRef: optionalString(result.ExternalRef),
	}); err != nil {
		return fmt.Errorf("mark release completed: %w", err)
	}

	return nil
}

func (p *StatePersistence) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	release, err := p.releaseRepository.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("find release before marking failed: %w", err)
	}
	if release == nil {
		return fmt.Errorf("release %s not found", id)
	}
	if release.Status != entity.ReleaseStatusRunning {
		return fmt.Errorf("release %s is not running", id)
	}

	status := entity.ReleaseStatusFailed
	if _, err := p.releaseRepository.UpdateOne(ctx, repository.UpdateReleaseInput{
		ID:     id,
		Status: &status,
	}); err != nil {
		return fmt.Errorf("mark release failed: %w", err)
	}

	_ = reason
	return nil
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
