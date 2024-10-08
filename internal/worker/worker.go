package worker

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/sirupsen/logrus"
)

// Worker processes tasks from a queue using registered task handlers
type Worker struct {
	id    int
	queue <-chan task.Task
}

// NewWorker creates and returns a new Worker instance with the given ID, task queue, and task handler registry
func NewWorker(id int, q <-chan task.Task) *Worker {
	return &Worker{
		id:    id,
		queue: q,
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
			result := t.Handler(t.Payload)
			if result.Error != nil {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_id":   t.ID,
					"error":     result.Error,
				}).Error("Failed to process task")
			}

			t.ResultCh <- result
		}
	}
}
