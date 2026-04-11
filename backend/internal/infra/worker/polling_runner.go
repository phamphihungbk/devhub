package worker

import core "devhub-backend/internal/infra/worker/core"

type PollingRunnerConfig = core.PollingRunnerConfig
type PollingRunner[T core.Identifiable, R any] = core.PollingRunner[T, R]

func NewPollingRunner[T core.Identifiable, R any](
	cfg PollingRunnerConfig,
	queue core.QueueSourceAdapter[T],
	persistence core.StatePersistence[R],
	executor core.Executor[T, R],
	observer core.Observability,
) (*PollingRunner[T, R], error) {
	return core.NewPollingRunner(cfg, queue, persistence, executor, observer)
}
