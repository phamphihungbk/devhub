package repository

import (
	"context"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type TeamRepository interface {
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Team, error)
}
