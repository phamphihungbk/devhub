package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	entity "github.com/phamphihungbk/devhub-backend/internal/domain/entity"
)

type UserRepository interface {
	CreateOne(ctx context.Context, user *entity.User) (*entity.User, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindAll(ctx context.Context, filter FindAllUsersFilter) (*entity.Users, int64, error)
	// Update(ctx context.Context, user *entity.User) error
	// Delete(ctx context.Context, id uuid.UUID) error
}

type FindAllUsersFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}
