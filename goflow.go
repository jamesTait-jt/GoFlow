package goflow

import (
	"context"
	"fmt"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/jamesTait-jt/GoFlow/internal/workerpool"
	"github.com/jamesTait-jt/GoFlow/pkg/store"
	"github.com/jamesTait-jt/GoFlow/worker"
)

type Broker interface {
	Submit(t task.Task)
	Dequeue() <-chan task.Task
}

type GoFlow struct {
	ctx          context.Context
	cancel       context.CancelFunc
	taskBroker   Broker
	workers      *workerpool.Pool
	taskHandlers store.KVStore[string, task.Handler]
	results      store.KVStore[string, task.Result]
}

func NewGoFlow(
	numWorkers int,
	workerFactory func(id int) worker.Worker,
	taskHandlers store.KVStore[string, task.Handler],
	results store.KVStore[string, task.Result],
	taskBroker Broker,
) *GoFlow {
	ctx, cancel := context.WithCancel(context.Background())

	workerPool := workerpool.NewWorkerPool(numWorkers, workerFactory)

	gf := GoFlow{
		ctx:          ctx,
		cancel:       cancel,
		taskBroker:   taskBroker,
		workers:      workerPool,
		taskHandlers: taskHandlers,
		results:      results,
	}

	return &gf
}

func (gf *GoFlow) Start() {
	gf.workers.Start(gf.ctx, gf.taskBroker)
}

func (gf *GoFlow) RegisterHandler(taskType string, handler task.Handler) {
	gf.taskHandlers.Put(taskType, handler)
}

func (gf *GoFlow) Push(taskType string, payload any) (string, error) {
	handler, ok := gf.taskHandlers.Get(taskType)
	if !ok {
		return "", fmt.Errorf("no handler defined for taskType: %s", taskType)
	}

	t := task.NewTask(taskType, payload, handler)

	gf.taskBroker.Submit(t)
	go gf.persistResult(t)

	return t.ID, nil
}

func (gf *GoFlow) GetResult(taskID string) (task.Result, bool) {
	result, ok := gf.results.Get(taskID)
	return result, ok
}

func (gf *GoFlow) Stop() {
	// Cancel the context, signalling to all the workers that they must stop
	gf.cancel()

	// Wait for all the workers to stop
	gf.workers.AwaitShutdown()
}

func (gf *GoFlow) persistResult(t task.Task) {
	result := <-t.ResultCh
	gf.results.Put(t.ID, result)
}
