package workerpool

import (
	"context"
	"sync"

	"github.com/google/uuid"
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

type Pool struct {
	workers map[string]Worker
	wg      *sync.WaitGroup
}

func NewWorkerPool(workers []Worker) *Pool {
	wp := &Pool{
		workers: make(map[string]Worker, len(workers)),
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < len(workers); i++ {
		id := uuid.New().String()
		wp.workers[id] = workers[i]
		wp.workers[id].SetID(id)
	}

	return wp
}

func (wp *Pool) Start(ctx context.Context, taskSource TaskSource, resultsCh chan<- task.Result) {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		worker.Start(ctx, wp.wg, taskSource, resultsCh)
	}
}

func (wp *Pool) AwaitShutdown() {
	wp.wg.Wait()
}
