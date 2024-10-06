package worker

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/sirupsen/logrus"
)

// taskHandlerRegistry defines an interface for retrieving task handlers based on task type
type taskHandlerRegistry interface {
	GetHandler(taskType string) (task.Handler, bool)
}

// Worker processes tasks from a queue using registered task handlers
type Worker struct {
	id                  int
	queue               <-chan task.Task
	taskHandlerRegistry taskHandlerRegistry
}

// NewWorker creates and returns a new Worker instance with the given ID, task queue, and task handler registry
func NewWorker(id int, q <-chan task.Task, taskHandlerRegistry taskHandlerRegistry) *Worker {
	return &Worker{
		id:                  id,
		queue:               q,
		taskHandlerRegistry: taskHandlerRegistry,
	}
}

// Start begins the worker's task processing in a separate goroutine.
// It takes a context to manage the worker's lifecycle and a WaitGroup to signal completion.
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.Infof("Worker %d starting...", w.id)

	go func() {
		defer wg.Done()
		w.processQueue(ctx)
	}()
}

// processQueue continuously listens for tasks from the queue and processes them.
// It will stop processing when the provided context is done.
func (w *Worker) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
			}).Info("Received shutdown signal, stopping worker")

			return

		case t := <-w.queue:
			handler, ok := w.taskHandlerRegistry.GetHandler(t.Type)
			if !ok {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_type": t.Type,
				}).Warn("No handler registered for task type")

				continue
			}

			err := handler(t.Payload)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_id":   t.ID,
					"error":     err,
				}).Error("Failed to process task")
			}
		}
	}
}
