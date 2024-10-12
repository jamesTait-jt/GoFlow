package workerpool

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jamesTait-jt/goflow/workerpool"
)

type Pool struct {
	workers map[string]workerpool.Worker
	wg      *sync.WaitGroup
}

func NewWorkerPool(workers []workerpool.Worker) *Pool {
	wp := &Pool{
		workers: make(map[string]workerpool.Worker, len(workers)),
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < len(workers); i++ {
		id := uuid.New().String()
		wp.workers[id] = workers[i]
		wp.workers[id].SetID(id)
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
