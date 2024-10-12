package workerpool

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/task"
)

type TaskSource interface {
	Dequeue() <-chan task.Task
}

type Worker interface {
	SetID(id string)
	Start(ctx context.Context, wg *sync.WaitGroup, taskSource TaskSource)
}
