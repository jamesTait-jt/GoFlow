package broker

import (
	"context"

	"github.com/jamesTait-jt/goflow/task"
)

// ChannelBroker is a task broker implementation that wraps a buffered Go channel.
// It manages task submission and retrieval for the GoFlow framework.
type ChannelBroker struct {
	taskQueue chan task.Task
}

// NewChannelBroker creates a new ChannelBroker with a buffered channel of the
// specified size. The buffer size determines how many tasks can be queued before
// further submissions block.
func NewChannelBroker(bufferSize int) *ChannelBroker {
	q := make(chan task.Task, bufferSize)

	return &ChannelBroker{taskQueue: q}
}

// Submit adds a task to the ChannelBroker's queue. If the queue is full, it will
// block until space is available.
func (cb *ChannelBroker) Submit(ctx context.Context, t task.Task) {
	select {
	// Without this, it is possible for this goroutine to be locked trying to
	// write to finished workers
	case <-ctx.Done():
		return

	case cb.taskQueue <- t:
		return
	}
}

// Dequeue returns a read-only channel of tasks, allowing workers to retrieve
// tasks for processing.
func (cb *ChannelBroker) Dequeue(_ context.Context) <-chan task.Task {
	return cb.taskQueue
}
