package workerpool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/jamesTait-jt/goflow/task"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTaskSource struct {
	mock.Mock
}

func (m *mockTaskSource) Dequeue(ctx context.Context) <-chan task.Task {
	args := m.Called(ctx)
	return args.Get(0).(chan task.Task)
}

func TestWorker_Start(t *testing.T) {
	t.Run("Processes a task successfully", func(t *testing.T) {
		// Arrange
		taskChan := make(chan task.Task)
		defer close(taskChan)

		ctx, cancel := context.WithCancel(context.Background())

		taskSource := mockTaskSource{}
		taskSource.On("Dequeue", ctx).Return(taskChan).Twice()

		var receivedPayload string

		expectedResult := task.Result{
			Payload: "payload",
			Error:   nil,
		}
		handler := task.Handler(func(payload any) task.Result {
			receivedPayload, _ = payload.(string)
			return expectedResult
		})

		resultCh := make(chan task.Result)
		taskToProcess := task.Task{
			ID:      "task-1",
			Type:    "test-task",
			Payload: "test-payload",
			Handler: handler,
		}
		expectedResult.TaskID = taskToProcess.ID

		w := LocalWorker{
			id: "1",
		}

		var wg sync.WaitGroup

		wg.Add(1)

		// Act
		w.Start(ctx, &wg, &taskSource, resultCh)
		taskChan <- taskToProcess

		result := <-resultCh

		cancel()

		// Assert
		wg.Wait()

		taskSource.AssertExpectations(t)

		assert.Equal(t, expectedResult, result)
		assert.Equal(t, taskToProcess.Payload, receivedPayload)
	})

	t.Run("Logs error when handler reports error", func(t *testing.T) {
		// Arrange
		var logOutput bytes.Buffer

		logrus.SetOutput(&logOutput)
		logrus.SetLevel(logrus.ErrorLevel)

		taskChan := make(chan task.Task)
		defer close(taskChan)

		ctx, cancel := context.WithCancel(context.Background())

		taskSource := mockTaskSource{}
		taskSource.On("Dequeue", ctx).Return(taskChan).Twice()

		var receivedPayload string

		expectedResult := task.Result{
			Payload: nil,
			Error:   errors.New("handler error"),
		}
		handler := task.Handler(func(payload any) task.Result {
			receivedPayload, _ = payload.(string)
			return expectedResult
		})

		resultCh := make(chan task.Result)
		taskToProcess := task.Task{
			ID:      "task-1",
			Type:    "test-task",
			Payload: "test-payload",
			Handler: handler,
		}
		expectedResult.TaskID = taskToProcess.ID

		w := LocalWorker{
			id: "1",
		}

		var wg sync.WaitGroup

		wg.Add(1)

		// Act
		w.Start(ctx, &wg, &taskSource, resultCh)
		taskChan <- taskToProcess

		result := <-resultCh

		cancel()

		// Assert
		wg.Wait()

		taskSource.AssertExpectations(t)

		assert.Equal(t, taskToProcess.Payload, receivedPayload)
		assert.Equal(t, expectedResult, result)

		require.Contains(t, logOutput.String(), "Failed to process task")
		require.Contains(t, logOutput.String(), "worker_id=1")
		require.Contains(t, logOutput.String(), "task_id=task-1")
		require.Contains(t, logOutput.String(), fmt.Sprintf("error=%q", expectedResult.Error))
	})
}
