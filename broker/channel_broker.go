package broker

import "github.com/jamesTait-jt/goflow/task"

type ChannelBroker struct {
	taskQueue chan task.Task
}

func NewChannelBroker(bufferSize int) *ChannelBroker {
	q := make(chan task.Task, bufferSize)

	return &ChannelBroker{taskQueue: q}
}

func (cb *ChannelBroker) Submit(t task.Task) {
	cb.taskQueue <- t
}

func (cb *ChannelBroker) Dequeue() <-chan task.Task {
	return cb.taskQueue
}
