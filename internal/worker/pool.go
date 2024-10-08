package worker

import (
	"sync"

	"github.com/jamesTait-jt/GoFlow/internal/task"
)

type Pool struct {
	wg      *sync.WaitGroup
	workers map[int]Worker
}

func NewPool(numWorkers int, taskQueue <-chan task.Task) *Pool {
	p := Pool{}

	for i := 0; i < numWorkers; i++ {
		p.workers[i] = NewWorker(i, taskQueue, gf.taskHandlers)
	}

	return &p
}
