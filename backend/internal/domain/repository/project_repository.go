package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateOne(ctx context.Context, project *entity.Project) (*entity.Project, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	FindAll(ctx context.Context, filter FindAllProjectsFilter) (*entity.Projects, int64, error)
	UpdateOne(ctx context.Context, input UpdateProjectInput) (*entity.Project, error)
	DeleteOne(ctx context.Context, id uuid.UUID) (*entity.Project, error)
}

type FindAllProjectsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}

type UpdateProjectInput struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	Environments *[]string
	Status       *entity.ProjectStatus
	OwnerTeam    *string
	RepoURL      *string
	RepoProvider *string
	OwnerContact *string
}
