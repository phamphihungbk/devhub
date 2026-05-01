package repository

import (
	"context"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type EnvironmentRepository interface {
	CreateOne(ctx context.Context, environment *entity.Environment) (*entity.Environment, error)
	FindAllByProjectID(ctx context.Context, projectID uuid.UUID) (*entity.Environments, error)
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}
