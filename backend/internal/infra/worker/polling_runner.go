package worker

import core "devhub-backend/internal/infra/worker/core"

type PollingRunnerConfig = core.PollingRunnerConfig
type PollingRunner[T core.Identifiable, R any] struct {
	*core.PollingRunner[T, R]
}

func NewPollingRunner[T core.Identifiable, R any](
	cfg PollingRunnerConfig,
	queue core.QueueSourceAdapter[T],
	persistence core.StatePersistence[R],
	executor core.Executor[T, R],
	observer core.Observability,
) (*PollingRunner[T, R], error) {
	runner, err := core.NewPollingRunner(cfg, queue, persistence, executor, observer)
	if err != nil {
		return nil, err
	}

	return &PollingRunner[T, R]{PollingRunner: runner}, nil
}
