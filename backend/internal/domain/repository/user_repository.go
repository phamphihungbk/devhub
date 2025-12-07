package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateOne(ctx context.Context, user *entity.User) (*entity.User, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindAll(ctx context.Context, filter FindAllUsersFilter) (*entity.Users, int64, error)
	UpdateOne(ctx context.Context, input UpdateUserInput) (*entity.User, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

type FindAllUsersFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdateUserInput struct {
	ID   uuid.UUID
	Name *string
	Role *entity.UserRole
}
