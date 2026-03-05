package runner

import (
	"context"

	"devhub-backend/internal/infra/worker/queue"
)

// Runner executes one claimed job.
type Runner interface {
	Run(ctx context.Context, job queue.Job) error
}
