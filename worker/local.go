package worker

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/goflow/task"
	"github.com/jamesTait-jt/goflow/workerpool"
	"github.com/sirupsen/logrus"
)

type LocalWorker struct {
	id string
}

func NewLocalWorker() *LocalWorker {
	return &LocalWorker{}
}

func (w *LocalWorker) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	taskSource workerpool.TaskSource,
	resultsCh chan<- task.Result,
) {
	logrus.Infof("Worker %s starting...", w.id)

	go func() {
		defer wg.Done()
		w.processQueue(ctx, taskSource, resultsCh)
	}()
}

func (w *LocalWorker) SetID(id string) {
	w.id = id
}

func (w *LocalWorker) processQueue(ctx context.Context, taskSource workerpool.TaskSource, resultsCh chan<- task.Result) {
	for {
		select {
		case <-ctx.Done():
			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
			}).Info("Received shutdown signal, stopping worker")

			return

		case t := <-taskSource.Dequeue(ctx):
			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
				"task_id":   t.ID,
			}).Info("Picked up task")

			result := t.Handler(t.Payload)
			result.TaskID = t.ID

			if result.Error != nil {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"task_id":   t.ID,
					"error":     result.Error,
				}).Error("Failed to process task")
			}

			resultsCh <- result
		}
	}
}
