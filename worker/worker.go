package worker

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/task"
)

type TaskSource interface {
	Dequeue() <-chan task.Task
}

type Worker interface {
	Start(ctx context.Context, wg *sync.WaitGroup, taskSource TaskSource)
}
