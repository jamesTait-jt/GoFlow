package worker

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type LocalWorker struct {
	id int
}

func NewWorker(id int) *LocalWorker {
	return &LocalWorker{
		id: id,
	}
}

func (w *LocalWorker) Start(ctx context.Context, wg *sync.WaitGroup, taskSource TaskSource) {
	logrus.Infof("Worker %d starting...", w.id)

	go func() {
		defer wg.Done()
		w.processQueue(ctx, taskSource)
	}()
}

func (w *LocalWorker) processQueue(ctx context.Context, taskSource TaskSource) {
	for {
		select {
		case <-ctx.Done():
			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
			}).Info("Received shutdown signal, stopping worker")

			return

		case t := <-taskSource.Dequeue():
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
