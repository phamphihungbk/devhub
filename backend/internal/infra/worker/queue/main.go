package queue

import (
	"context"
	"errors"
)

var ErrNoJobAvailable = errors.New("no job available")

// Job is transport-agnostic so DB/Kafka/SQS backends can map into it.
type Job struct {
	ID       string
	Payload  []byte
	Attempts int
}

// Queue defines the contract for job sources.
type Queue interface {
	ClaimNextPending(ctx context.Context) (*Job, error)
	MarkSucceeded(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, failureReason string) error
}
