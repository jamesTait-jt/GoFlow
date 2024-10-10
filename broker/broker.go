package broker

import "github.com/jamesTait-jt/GoFlow/internal/task"

type Broker interface {
	Submit(t task.Task)
	Dequeue() <-chan task.Task
}
