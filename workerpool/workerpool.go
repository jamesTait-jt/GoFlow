package workerpool

import (
	"context"
	"sync"

	"github.com/jamesTait-jt/goflow/task"
	"github.com/sirupsen/logrus"
)

type HandlerGetter interface {
	Get(taskType string) (task.Handler, bool)
}

type Pool struct {
	numWorkers int
	wg         *sync.WaitGroup
}

func New(numWorkers int) *Pool {
	wp := &Pool{
		numWorkers: numWorkers,
		wg:         &sync.WaitGroup{},
	}

	return wp
}

func (wp *Pool) Start(
	ctx context.Context,
	taskQueue task.Dequeuer[task.Task],
	results task.Submitter[task.Result],
	taskHandlers HandlerGetter,
) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go worker(ctx, wp.wg, taskQueue, results, taskHandlers)
	}
}

func (wp *Pool) AwaitShutdown() {
	wp.wg.Wait()
}

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	taskQueue task.Dequeuer[task.Task],
	results task.Submitter[task.Result],
	taskHandlers HandlerGetter,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logrus.Info("Received shutdown signal, stopping worker")
			return

		case t := <-taskQueue.Dequeue(ctx):
			logrus.WithFields(logrus.Fields{
				"task_id": t.ID,
			}).Info("Picked up task")

			handler, ok := taskHandlers.Get(t.Type)
			if !ok {
				logrus.WithFields(logrus.Fields{
					"task_type": t.Type,
				}).Error("No handler registered for task type")

				continue
			}

			result := handler(t.Payload)
			result.TaskID = t.ID

			if result.Error != nil {
				logrus.WithFields(logrus.Fields{
					"task_id": t.ID,
					"error":   result.Error,
				}).Error("Failed to process task")
			}

			results.Submit(ctx, result)
		}
	}
}
