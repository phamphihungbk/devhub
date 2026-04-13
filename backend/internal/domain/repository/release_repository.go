package repository

import (
	"context"

	entity "devhub-backend/internal/domain/entity"
)

type ReleaseRepository interface {
	CreateOne(ctx context.Context, release *entity.Release) (*entity.Release, error)
}
