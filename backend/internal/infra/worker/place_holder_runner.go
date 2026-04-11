package worker

import (
	"context"
	"time"

	infraLogger "devhub-backend/internal/infra/logger"
)

type placeholderRunner struct {
	name      string
	pollDelay time.Duration
	logger    infraLogger.Logger
}

func newPlaceholderRunner(name string, pollDelay time.Duration, logger infraLogger.Logger) *placeholderRunner {
	return &placeholderRunner{
		name:      name,
		pollDelay: pollDelay,
		logger:    logger,
	}
}

func (r *placeholderRunner) Name() string {
	return r.name
}

func (r *placeholderRunner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.pollDelay)
	defer ticker.Stop()

	r.logger.Warn(ctx, "runner is in placeholder mode; wire queue/executor to enable processing", infraLogger.Fields{
		"runner": r.name,
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
