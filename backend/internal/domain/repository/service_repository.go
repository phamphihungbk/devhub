package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type ServiceRepository interface {
	CreateOne(ctx context.Context, service *entity.Service) (*entity.Service, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	FindAll(ctx context.Context, filter FindAllServicesFilter) (*entity.Services, int64, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.Service, error)
}

type FindAllServicesFilter struct {
	ProjectID uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}
