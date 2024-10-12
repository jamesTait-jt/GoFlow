package workerpool

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/workerpool"
)

type Pool struct {
	workers map[int]workerpool.Worker
	wg      *sync.WaitGroup
}

func NewWorkerPool(numWorkers int, workerFactory func(id int) workerpool.Worker) *Pool {
	wp := &Pool{
		workers: make(map[int]workerpool.Worker, numWorkers),
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < numWorkers; i++ {
		wp.workers[i] = workerFactory(i)
	}

	return wp
}

func (wp *Pool) Start(ctx context.Context, taskSource workerpool.TaskSource) {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		worker.Start(ctx, wp.wg, taskSource)
	}
}

func (wp *Pool) AwaitShutdown() {
	wp.wg.Wait()
}
