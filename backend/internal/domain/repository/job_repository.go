package repository

import (
	"context"
	"time"

	entity "devhub-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type JobRepository interface {
	CreateOne(ctx context.Context, job *entity.Job) (*entity.Job, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Job, error)
	FindOnePending(ctx context.Context, jobType entity.JobType) (*entity.Job, error)
	UpdateOne(ctx context.Context, input UpdateJobInput) (*entity.Job, error)
}

type UpdateJobInput struct {
	ID         uuid.UUID
	Status     *entity.JobStatus
	Result     *string
	Error      *string
	Attempts   *int
	StartedAt  *time.Time
	FinishedAt *time.Time
}
