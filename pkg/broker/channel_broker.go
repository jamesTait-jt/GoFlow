package broker

import (
	"context"

	"github.com/jamesTait-jt/goflow/pkg/task"
)

// ChannelBroker is a task broker implementation that wraps a buffered Go channel.
// It manages task submission and retrieval for the GoFlow framework.
type ChannelBroker[T task.TaskOrResult] struct {
	taskQueue chan T
}

// NewChannelBroker creates a new ChannelBroker with a buffered channel of the
// specified size. The buffer size determines how many tasks can be queued before
// further submissions block.
func NewChannelBroker[T task.TaskOrResult](bufferSize int) *ChannelBroker[T] {
	q := make(chan T, bufferSize)

	return &ChannelBroker[T]{taskQueue: q}
}

// Submit adds a task to the ChannelBroker's queue. If the queue is full, it will
// block until space is available.
func (cb *ChannelBroker[T]) Submit(ctx context.Context, t T) error {
	select {
	// Without this, it is possible for this goroutine to be locked trying to
	// write to finished workers
	case <-ctx.Done():
		return nil

	case cb.taskQueue <- t:
		return nil
	}
}

// Dequeue returns a read-only channel of tasks, allowing workers to retrieve
// tasks for processing.
func (cb *ChannelBroker[T]) Dequeue(_ context.Context) <-chan T {
	return cb.taskQueue
}
