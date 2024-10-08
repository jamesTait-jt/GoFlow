package worker

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/internal/task"
)

type Pool struct {
	workers map[int]*Worker
	ctx     context.Context
	wg      *sync.WaitGroup
}

func NewWorkerPool(numWorkers int, queue <-chan task.Task, ctx context.Context) *Pool {
	wp := &Pool{
		workers: make(map[int]*Worker, numWorkers),
		ctx:     ctx,
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < numWorkers; i++ {
		wp.workers[i] = NewWorker(i, queue)
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