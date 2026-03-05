package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"devhub-backend/internal/infra/worker/queue"
	"devhub-backend/internal/infra/worker/runner"
	"devhub-backend/internal/infra/logger"
)

type Worker struct {
	cfg    Config
	queue  queue.Queue
	runner runner.Runner
	logger logger.Logger
}

func New(cfg Config, q queue.Queue, r runner.Runner, l Logger) *Worker {
	cfg = cfg.withDefaults()

	return &Worker{
		cfg:    cfg,
		queue:  q,
		runner: r,
		logger: l,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info(ctx, "worker started", map[string]any{
		"poll_interval": w.cfg.PollInterval.String(),
		"max_workers":   w.cfg.MaxWorkers,
	})

	var wg sync.WaitGroup
	for i := 0; i < w.cfg.MaxWorkers; i++ {
		workerID := i + 1
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			w.runLoop(ctx, id)
		}(workerID)
	}

	<-ctx.Done()
	wg.Wait()
	return nil
}

func (w *Worker) runLoop(ctx context.Context, workerID int) {
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.tick(ctx, workerID); err != nil {
				w.logger.Error(ctx, "worker tick failed", err, map[string]any{
					"worker_id": workerID,
				})
			}
		}
	}
}

func (w *Worker) tick(ctx context.Context, workerID int) error {
	job, err := w.queue.ClaimNextPending(ctx)
	if err != nil {
		if errors.Is(err, queue.ErrNoJobAvailable) {
			return nil
		}
		return fmt.Errorf("claim next pending job: %w", err)
	}

	if w.cfg.RetryPolicy.MaxAttempts > 0 && job.Attempts > w.cfg.RetryPolicy.MaxAttempts {
		retryErr := fmt.Sprintf("max attempts exceeded: %d", w.cfg.RetryPolicy.MaxAttempts)
		if err := w.queue.MarkFailed(ctx, job.ID, retryErr); err != nil {
			return fmt.Errorf("mark failed after max attempts for job %s: %w", job.ID, err)
		}
		w.logger.Warn(ctx, "job exceeded max attempts", map[string]any{
			"worker_id": workerID,
			"job_id":    job.ID,
			"attempts":  job.Attempts,
		})
		return nil
	}

	if err := w.runner.Run(ctx, *job); err != nil {
		if markErr := w.queue.MarkFailed(ctx, job.ID, err.Error()); markErr != nil {
			return fmt.Errorf("run error: %v; mark failed error: %w", err, markErr)
		}
		return fmt.Errorf("run job %s: %w", job.ID, err)
	}

	if err := w.queue.MarkSucceeded(ctx, job.ID); err != nil {
		return fmt.Errorf("mark succeeded for job %s: %w", job.ID, err)
	}

	return nil
}
