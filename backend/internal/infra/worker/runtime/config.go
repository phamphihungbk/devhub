package worker

import "time"

const (
	defaultPollInterval = 2 * time.Second
	defaultMaxWorkers   = 1
)

type RetryPolicy struct {
	MaxAttempts int
}

type Config struct {
	PollInterval time.Duration
	MaxWorkers   int
	RetryPolicy  RetryPolicy
}

func (c Config) withDefaults() Config {
	if c.PollInterval <= 0 {
		c.PollInterval = defaultPollInterval
	}

	if c.MaxWorkers <= 0 {
		c.MaxWorkers = defaultMaxWorkers
	}

	return c
}
