package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type DeploymentRepository interface {
	CreateOne(ctx context.Context, project *entity.Project) (*entity.Deployment, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Deployment, error)
	FindAll(ctx context.Context, filter FindAllDeploymentsFilter) (*entity.Deployments, int64, error)
	UpdateOne(ctx context.Context, input UpdateDeploymentInput) (*entity.Deployment, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.Deployment, error)
}

type FindAllDeploymentsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdateDeploymentInput struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	Environments *[]string
}
