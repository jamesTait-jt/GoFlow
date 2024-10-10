package broker

import "github.com/jamesTait-jt/GoFlow/task"

type Broker interface {
	Submit(t task.Task)
	Dequeue() <-chan task.Task
}
