package repository

import (
	"context"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type TeamRepository interface {
	CreateOne(ctx context.Context, team *entity.Team) (*entity.Team, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	FindAll(ctx context.Context, filter FindAllTeamsFilter) (*entity.Teams, int64, error)
	UpdateOne(ctx context.Context, input UpdateTeamInput) (*entity.Team, error)
}

type FindAllTeamsFilter struct {
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdateTeamInput struct {
	ID           uuid.UUID
	Name         *string
	OwnerContact *string
}
