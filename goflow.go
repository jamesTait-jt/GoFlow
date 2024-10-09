package goflow

import (
	"context"
	"fmt"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/jamesTait-jt/GoFlow/internal/worker"
)

type TaskHandlerRegistry interface {
	RegisterHandler(taskType string, handler task.Handler)
	GetHandler(taskType string) (task.Handler, bool)
}

type Broker interface {
	Submit(t task.Task)
	Dequeue() <-chan task.Task
}

type GoFlow struct {
	workers      *worker.Pool
	taskBroker   Broker
	taskHandlers TaskHandlerRegistry
	results      map[string]task.Result
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewGoFlow(numWorkers int, taskBroker Broker, taskHandlerRegistry TaskHandlerRegistry) *GoFlow {
	ctx, cancel := context.WithCancel(context.Background())

	workerPool := worker.NewWorkerPool(ctx, numWorkers, taskBroker)

	gf := GoFlow{
		workers:      workerPool,
		taskBroker:   taskBroker,
		taskHandlers: taskHandlerRegistry,
		ctx:          ctx,
		cancel:       cancel,
	}

	return &gf
}

func (gf *GoFlow) Start() {
	gf.workers.Start()
}

func (gf *GoFlow) RegisterHandler(taskType string, handler task.Handler) {
	gf.taskHandlers.RegisterHandler(taskType, handler)
}

func (gf *GoFlow) Push(taskType string, payload any) (string, error) {
	handler, ok := gf.taskHandlers.GetHandler(taskType)
	if !ok {
		return "", fmt.Errorf("no handler defined for taskType: %s", taskType)
	}

	t := task.NewTask(taskType, payload, handler)

	gf.taskBroker.Submit(t)

	go gf.persistResult(t)

	return t.ID, nil
}

func (gf *GoFlow) GetResult(taskID string) (task.Result, bool) {
	result, ok := gf.results[taskID]
	return result, ok
}

func (gf *GoFlow) Stop() {
	// Cancel the context, signalling to all the workers that they must stop
	gf.cancel()

	// Wait for all the workers to stop
	gf.workers.WaitForShutdown()
}

func (gf *GoFlow) persistResult(t task.Task) {
	result := <-t.ResultCh
	gf.results[t.ID] = result
}
