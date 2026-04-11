package core

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

const DefaultPollDelay = 3 * time.Second

type Runner interface {
	Name() string
	Run(ctx context.Context) error
}

type Identifiable interface {
	GetID() uuid.UUID
}

type QueueSourceAdapter[T Identifiable] interface {
	Dequeue(ctx context.Context) (*T, error)
}

type StatePersistence[R any] interface {
	MarkRunning(ctx context.Context, id uuid.UUID) error
	MarkCompleted(ctx context.Context, id uuid.UUID, result R) error
	MarkFailed(ctx context.Context, id uuid.UUID, reason string) error
}

type Executor[T any, R any] interface {
	Execute(ctx context.Context, job *T) (R, error)
}

type Observability interface {
	OnRunnerStart(ctx context.Context, runner string, pollDelay time.Duration)
	OnRunnerStop(ctx context.Context, runner string, reason error)
	OnJobDequeued(ctx context.Context, runner string, jobID uuid.UUID)
	OnJobCompleted(ctx context.Context, runner string, jobID uuid.UUID)
	OnError(ctx context.Context, runner string, phase string, err error, jobID uuid.UUID)
}

type NoopObservability struct{}

func (NoopObservability) OnRunnerStart(context.Context, string, time.Duration) {}
func (NoopObservability) OnRunnerStop(context.Context, string, error)          {}
func (NoopObservability) OnJobDequeued(context.Context, string, uuid.UUID)     {}
func (NoopObservability) OnJobCompleted(context.Context, string, uuid.UUID)    {}
func (NoopObservability) OnError(context.Context, string, string, error, uuid.UUID) {
}

type PollingRunnerConfig struct {
	Name      string
	PollDelay time.Duration
}

type PollingRunner[T Identifiable, R any] struct {
	name        string
	pollDelay   time.Duration
	queue       QueueSourceAdapter[T]
	persistence StatePersistence[R]
	executor    Executor[T, R]
	observer    Observability
}

func NewPollingRunner[T Identifiable, R any](
	cfg PollingRunnerConfig,
	queue QueueSourceAdapter[T],
	persistence StatePersistence[R],
	executor Executor[T, R],
	observer Observability,
) (*PollingRunner[T, R], error) {
	if cfg.Name == "" {
		return nil, errors.New("runner name is required")
	}

	if cfg.PollDelay <= 0 {
		cfg.PollDelay = DefaultPollDelay
	}

	if queue == nil {
		return nil, errors.New("queue source adapter is required")
	}

	if persistence == nil {
		return nil, errors.New("state persistence is required")
	}

	if executor == nil {
		return nil, errors.New("executor is required")
	}

	if observer == nil {
		observer = NoopObservability{}
	}

	return &PollingRunner[T, R]{
		name:        cfg.Name,
		pollDelay:   cfg.PollDelay,
		queue:       queue,
		persistence: persistence,
		executor:    executor,
		observer:    observer,
	}, nil
}

func (r *PollingRunner[T, R]) Name() string {
	return r.name
}

func (r *PollingRunner[T, R]) Run(ctx context.Context) error {
	r.observer.OnRunnerStart(ctx, r.name, r.pollDelay)

	ticker := time.NewTicker(r.pollDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.observer.OnRunnerStop(ctx, r.name, ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			r.processNext(ctx)
		}
	}
}

func (r *PollingRunner[T, R]) processNext(ctx context.Context) {
	job, err := r.queue.Dequeue(ctx)
	if err != nil {
		r.observer.OnError(ctx, r.name, "dequeue", err, uuid.Nil)
		return
	}

	if job == nil {
		return
	}

	jobID := (*job).GetID()
	r.observer.OnJobDequeued(ctx, r.name, jobID)

	if err := r.persistence.MarkRunning(ctx, jobID); err != nil {
		r.observer.OnError(ctx, r.name, "mark_running", err, jobID)
		return
	}

	result, err := r.executor.Execute(ctx, job)
	if err != nil {
		r.observer.OnError(ctx, r.name, "execute", err, jobID)
		if markErr := r.persistence.MarkFailed(ctx, jobID, err.Error()); markErr != nil {
			r.observer.OnError(ctx, r.name, "mark_failed", markErr, jobID)
		}
		return
	}

	if err := r.persistence.MarkCompleted(ctx, jobID, result); err != nil {
		r.observer.OnError(ctx, r.name, "mark_completed", err, jobID)
		return
	}

	r.observer.OnJobCompleted(ctx, r.name, jobID)
}
