package repository

import (
	"context"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type ReleaseRepository interface {
	CreateOne(ctx context.Context, release *entity.Release) (*entity.Release, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Release, error)
	FindOnePending(ctx context.Context) (*entity.Release, error)
	UpdateOne(ctx context.Context, input UpdateReleaseInput) (*entity.Release, error)
}

type UpdateReleaseInput struct {
	ID          uuid.UUID
	Status      *entity.ReleaseStatus
	ExternalRef *string
}
