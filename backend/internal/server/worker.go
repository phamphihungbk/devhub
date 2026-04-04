package server

import (
	"context"
	"devhub-backend/internal/config"
	"time"

	infraWorker "devhub-backend/internal/infra/worker"
)

type Worker struct {
	cfg         *config.Config
	concurrency int
	pollDelay   time.Duration
	workerTypes []string
}

type WorkerOption func(*Worker)

func WithWorkerConcurrency(n int) WorkerOption {
	return func(w *Worker) {
		w.concurrency = n
	}
}

func WithWorkerPollInterval(d time.Duration) WorkerOption {
	return func(w *Worker) {
		w.pollDelay = d
	}
}

func WithWorkerTypes(types []string) WorkerOption {
	return func(w *Worker) {
		w.workerTypes = types
	}
}

func NewWorker(opts ...WorkerOption) *Worker {
	w := &Worker{
		cfg:         config.MustConfigure(),
		concurrency: 1,
		pollDelay:   3 * time.Second,
		workerTypes: []string{"scaffold", "deployment"},
	}

	for _, opt := range opts {
		opt(w)
	}

	if w.concurrency <= 0 {
		w.concurrency = 1
	}
	if w.pollDelay <= 0 {
		w.pollDelay = 3 * time.Second
	}

	return w
}

func (w *Worker) Start() error {
	ctx := context.Background()

	runtime, err := infraWorker.Bootstrap(w.cfg)

	if err != nil {
		return err
	}

	runners, err := infraWorker.BuildRunnersWithConfig(runtime.Logger(), runtime.DB(), infraWorker.BuildRunnersConfig{
		WorkerTypes: w.workerTypes,
		PollDelay:   w.pollDelay,
	})

	if err != nil {
		_ = runtime.DB().Close()
		return err
	}

	return infraWorker.Run(ctx, runtime, runners, infraWorker.RunConfig{
		Concurrency:     w.concurrency,
		PollDelay:       w.pollDelay,
		ShutdownTimeout: 30 * time.Second,
	})
}
