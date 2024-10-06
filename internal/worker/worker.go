package worker

import (
	"context"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/sirupsen/logrus"
)

type taskHandlerRegistry interface {
	GetHandler(taskType string) (task.TaskHandler, bool)
}

type Worker struct {
	id           int
	queue        <-chan task.Task
	taskHandlerRegistry taskHandlerRegistry
}

func NewWorker(id int, q <-chan task.Task, taskHandlerRegistry taskHandlerRegistry) *Worker {
	return &Worker{
		id:    id,
		queue: q,
		taskHandlerRegistry: taskHandlerRegistry,
	}
}

func (w *Worker) Start(ctx context.Context) {
	logrus.Infof("Worker %d starting...", w.id)

	go w.processQueue(ctx)
}

func (w *Worker) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
			}).Info("Received shutdown signal, stopping worker")

			return

		case task := <-w.queue:

			handler, ok := w.taskHandlerRegistry.GetHandler(task.Type)
			if !ok {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_type": task.Type,
				}).Warn("No handler registered for task type")

				continue
			}

			err := handler(task.Payload)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_id":   task.ID,
					"error":     err,
				}).Error("Failed to process task")
			}
		}
	}
}
