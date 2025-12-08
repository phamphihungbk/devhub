package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type ScaffoldRequestRepository interface {
	CreateOne(ctx context.Context, scaffoldRequest *entity.ScaffoldRequest) (*entity.ScaffoldRequest, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.ScaffoldRequest, error)
	FindAll(ctx context.Context, filter FindAllScaffoldRequestsFilter) (*entity.ScaffoldRequests, int64, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.ScaffoldRequest, error)
}

type FindAllScaffoldRequestsFilter struct {
	ProjectID uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdateScaffoldRequestInput struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	Environments *[]string
}
