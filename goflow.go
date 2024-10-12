package goflow

import (
	"context"
	"fmt"

	"github.com/jamesTait-jt/goflow/internal/workerpool"
	"github.com/jamesTait-jt/goflow/task"
	publicWorkerpool "github.com/jamesTait-jt/goflow/workerpool"
)

type Broker interface {
	Submit(t task.Task)
	Dequeue() <-chan task.Task
}

type workerPool interface {
	Start(ctx context.Context, taskSource publicWorkerpool.TaskSource)
	AwaitShutdown()
}

type KVStore[K comparable, V any] interface {
	Put(k K, v V)
	Get(k K) (V, bool)
}

type goflow struct {
	ctx          context.Context
	cancel       context.CancelFunc
	taskBroker   Broker
	workers      workerPool
	taskHandlers KVStore[string, task.Handler]
	results      KVStore[string, task.Result]
}

func Newgoflow(
	workers []publicWorkerpool.Worker,
	taskHandlerStore KVStore[string, task.Handler],
	resultsStore KVStore[string, task.Result],
	taskBroker Broker,
) *goflow {
	ctx, cancel := context.WithCancel(context.Background())

	workerPool := workerpool.NewWorkerPool(workers)

	gf := goflow{
		ctx:          ctx,
		cancel:       cancel,
		taskBroker:   taskBroker,
		workers:      workerPool,
		taskHandlers: taskHandlerStore,
		results:      resultsStore,
	}

	return &gf
}

func (gf *goflow) Start() {
	gf.workers.Start(gf.ctx, gf.taskBroker)
}

func (gf *goflow) RegisterHandler(taskType string, handler task.Handler) {
	gf.taskHandlers.Put(taskType, handler)
}

func (gf *goflow) Push(taskType string, payload any) (string, error) {
	handler, ok := gf.taskHandlers.Get(taskType)
	if !ok {
		return "", fmt.Errorf("no handler defined for taskType: %s", taskType)
	}

	t := task.NewTask(taskType, payload, handler)

	gf.taskBroker.Submit(t)
	go gf.persistResult(t)

	return t.ID, nil
}

func (gf *goflow) GetResult(taskID string) (task.Result, bool) {
	result, ok := gf.results.Get(taskID)
	return result, ok
}

func (gf *goflow) Stop() {
	// Cancel the context, signalling to all the workers that they must stop
	gf.cancel()

	// Wait for all the workers to stop
	gf.workers.AwaitShutdown()
}

func (gf *goflow) persistResult(t task.Task) {
	result := <-t.ResultCh
	gf.results.Put(t.ID, result)
}
