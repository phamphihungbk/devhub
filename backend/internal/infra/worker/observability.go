package worker

import (
	"context"
	"time"

	infraLogger "devhub-backend/internal/infra/logger"
	core "devhub-backend/internal/infra/worker/core"

	"github.com/google/uuid"
)

type Observability = core.Observability
type NoopObservability = core.NoopObservability

type LoggerObservability struct {
	logger infraLogger.Logger
}

func NewLoggerObservability(logger infraLogger.Logger) Observability {
	return LoggerObservability{logger: logger}
}

func (o LoggerObservability) OnRunnerStart(ctx context.Context, runner string, pollDelay time.Duration) {
	o.logger.Info(ctx, "polling runner started", infraLogger.Fields{
		"runner":     runner,
		"poll_delay": pollDelay.String(),
	})
}

func (o LoggerObservability) OnRunnerStop(ctx context.Context, runner string, reason error) {
	fields := infraLogger.Fields{"runner": runner}
	if reason != nil {
		fields["reason"] = reason.Error()
	}
	o.logger.Info(ctx, "polling runner stopped", fields)
}

func (o LoggerObservability) OnJobDequeued(ctx context.Context, runner string, jobID uuid.UUID) {
	o.logger.Info(ctx, "job dequeued", infraLogger.Fields{
		"runner": runner,
		"job_id": jobID.String(),
	})
}

func (o LoggerObservability) OnJobCompleted(ctx context.Context, runner string, jobID uuid.UUID) {
	o.logger.Info(ctx, "job completed", infraLogger.Fields{
		"runner": runner,
		"job_id": jobID.String(),
	})
}

func (o LoggerObservability) OnError(ctx context.Context, runner string, phase string, err error, jobID uuid.UUID) {
	fields := infraLogger.Fields{
		"runner": runner,
		"phase":  phase,
	}
	if jobID != uuid.Nil {
		fields["job_id"] = jobID.String()
	}
	o.logger.Error(ctx, "worker runner error", err, fields)
}
