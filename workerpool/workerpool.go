package workerpool

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/goflow/task"
)

type TaskSource interface {
	Dequeue(ctx context.Context) <-chan task.Task
}

type Worker interface {
	SetID(id string)
	Start(
		ctx context.Context,
		wg *sync.WaitGroup,
		taskSource TaskSource,
		resultsCh chan<- task.Result,
	)
}
