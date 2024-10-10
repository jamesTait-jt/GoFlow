package workerpool

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/worker"
)

type Pool struct {
	workers map[int]worker.Worker
	wg      *sync.WaitGroup
}

func NewWorkerPool(numWorkers int, workerFactory func(id int) worker.Worker) *Pool {
	wp := &Pool{
		workers: make(map[int]worker.Worker, numWorkers),
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < numWorkers; i++ {
		wp.workers[i] = workerFactory(i)
	}

	return wp
}

func (wp *Pool) Start(ctx context.Context, taskSource worker.TaskSource) {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		worker.Start(ctx, wp.wg, taskSource)
	}
}

func (wp *Pool) AwaitShutdown() {
	wp.wg.Wait()
}
