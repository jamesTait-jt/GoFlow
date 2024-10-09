package worker

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	id    int
	tasks taskSource
}

func NewWorker(id int, taskSource taskSource) *Worker {
	return &Worker{
		id:    id,
		tasks: taskSource,
	}
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.Infof("Worker %d starting...", w.id)

	go func() {
		defer wg.Done()
		w.processQueue(ctx)
	}()
}

func (w *Worker) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
			}).Info("Received shutdown signal, stopping worker")

			return

		case t := <-w.tasks.Dequeue():
			result := t.Handler(t.Payload)
			if result.Error != nil {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_id":   t.ID,
					"error":     result.Error,
				}).Error("Failed to process task")
			}

			t.ResultCh <- result
			close(t.ResultCh)
		}
	}
}
