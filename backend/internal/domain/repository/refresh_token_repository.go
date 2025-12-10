package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	CreateOne(ctx context.Context, input CreateRefreshTokenInput) (*entity.RefreshToken, error)
	FindOne(ctx context.Context, input FindOneRefreshTokenInput) (*entity.RefreshToken, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error)
}

type CreateRefreshTokenInput struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type FindOneRefreshTokenInput struct {
	UserID uuid.UUID
}
