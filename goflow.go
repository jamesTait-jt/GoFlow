package goflow

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/jamesTait-jt/GoFlow/internal/worker"
)

var QueueBufferSize = 100

type HandlerRegistry interface {
	RegisterHandler(taskType string, handler task.Handler)
	GetHandler(taskType string) (task.Handler, bool)
}

type GoFlow struct {
	workers      map[int]*worker.Worker
	taskQueue    chan task.Task
	taskHandlers HandlerRegistry
	ctx          context.Context
}

func NewGoFlow(numWorkers int) *GoFlow {
	gf := GoFlow{
		taskQueue:    make(chan task.Task, QueueBufferSize),
		taskHandlers: task.NewHandlerRegistry(),
	}

	return &gf
}

func (gf *GoFlow) Start() {
	var wg *sync.WaitGroup
	for _, worker := range gf.workers {
		wg.Add(1)
		worker.Start(gf.ctx, wg)
	}
}

func (gf *GoFlow) Stop() {

}
