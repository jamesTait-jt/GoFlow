package worker

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/internal/task"
)

type taskSource interface {
	Dequeue() <-chan task.Task
}

type worker interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type Pool struct {
	workers map[int]worker
	ctx     context.Context
	wg      *sync.WaitGroup
}

func NewWorkerPool(ctx context.Context, numWorkers int, taskSource taskSource) *Pool {
	wp := &Pool{
		workers: make(map[int]worker, numWorkers),
		ctx:     ctx,
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < numWorkers; i++ {
		wp.workers[i] = NewWorker(i, taskSource)
	}

	return wp
}

func (wp *Pool) Start() {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		worker.Start(wp.ctx, wp.wg)
	}
}

func (wp *Pool) WaitForShutdown() {
	wp.wg.Wait()
}
