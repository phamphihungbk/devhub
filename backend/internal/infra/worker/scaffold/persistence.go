package scaffold

import (
	"context"
	"fmt"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

type ScaffoldStatePersistence struct {
	scaffoldRequestRepository repository.ScaffoldRequestRepository
	serviceRepository         repository.ServiceRepository
}

var _ core.StatePersistence[ScaffoldExecutionResult] = (*ScaffoldStatePersistence)(nil)

func NewScaffoldStatePersistence(
	scaffoldRequestRepository repository.ScaffoldRequestRepository,
	serviceRepository repository.ServiceRepository,
) *ScaffoldStatePersistence {
	return &ScaffoldStatePersistence{scaffoldRequestRepository: scaffoldRequestRepository, serviceRepository: serviceRepository}
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

	if result.ProjectID == uuid.Nil {
		return fmt.Errorf("project id is required")
	}

	serviceName := strings.TrimSpace(result.ServiceName)
	if serviceName == "" {
		return fmt.Errorf("service name is required")
	}

	repoURL := strings.TrimSpace(result.RepoURL)
	if repoURL == "" {
		return fmt.Errorf("repo url is required")
	}

	if _, err := p.serviceRepository.CreateOne(ctx, &entity.Service{
		ProjectID: result.ProjectID,
		Name:      serviceName,
		RepoURL:   repoURL,
	}); err != nil {
		return err
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
