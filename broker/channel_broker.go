package broker

import "github.com/jamesTait-jt/goflow/task"

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
func (cb *ChannelBroker) Submit(t task.Task) {
	cb.taskQueue <- t
}

// Dequeue returns a read-only channel of tasks, allowing workers to retrieve
// tasks for processing.
func (cb *ChannelBroker) Dequeue() <-chan task.Task {
	return cb.taskQueue
}
