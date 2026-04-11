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
	FindOnePending(ctx context.Context) (*entity.ScaffoldRequest, error)
	FindAll(ctx context.Context, filter FindAllScaffoldRequestsFilter) (*entity.ScaffoldRequests, int64, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.ScaffoldRequest, error)
	UpdateOne(ctx context.Context, input UpdateScaffoldRequestInput) (*entity.ScaffoldRequest, error)
}

type FindAllScaffoldRequestsFilter struct {
	ProjectID uuid.UUID
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdateScaffoldRequestInput struct {
	ID            uuid.UUID
	PluginID      *uuid.UUID
	ProjectID     *uuid.UUID
	RequestedBy   *uuid.UUID
	Template      *string
	Status        *entity.ScaffoldRequestStatus
	Environment   *entity.ProjectEnvironment
	Variables     *entity.ScaffoldRequestVariables
	ApprovedBy    *uuid.UUID
	ResultRepoURL *string
	ApprovedAt    *time.Time
}
