package goflow

import (
	"context"
	"fmt"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/jamesTait-jt/GoFlow/internal/worker"
)

var queueBufferSize = 100

type HandlerRegistry interface {
	RegisterHandler(taskType string, handler task.Handler)
	GetHandler(taskType string) (task.Handler, bool)
}

type Dispatcher interface {
	Dispatch(t task.Task)
}

type GoFlow struct {
	workers        *worker.Pool
	taskQueue      chan task.Task
	taskDispatcher Dispatcher
	taskHandlers   HandlerRegistry
	results        map[string]task.Result
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewGoFlow(numWorkers int) *GoFlow {
	ctx, cancel := context.WithCancel(context.Background())

	taskQueue := make(chan task.Task, queueBufferSize)
	taskHandlers := task.NewHandlerRegistry()

	workerPool := worker.NewWorkerPool(numWorkers, taskQueue, ctx)

	gf := GoFlow{
		workers:      workerPool,
		taskQueue:    taskQueue,
		taskHandlers: taskHandlers,
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

	gf.taskQueue <- t

	go func() {
		result := <-t.ResultCh
		gf.results[t.ID] = result
	}()

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
