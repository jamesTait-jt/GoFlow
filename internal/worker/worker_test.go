package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockTaskHandlerRegistry is a mock for the task registry interface
type mockTaskHandlerRegistry struct {
	mock.Mock
}

func (m *mockTaskHandlerRegistry) GetHandler(taskType string) (task.Handler, bool) {
	args := m.Called(taskType)
	if args.Get(0) == nil {
		return nil, args.Bool(1) // Return nil as the task handler
	}

	return args.Get(0).(task.Handler), args.Bool(1)
}

func TestNewWorker_Creates_a_new_worker_with_variables_initialised(t *testing.T) {
	// Arrange
	queue := make(<-chan task.Task)
	registry := new(mockTaskHandlerRegistry)

	// Act
	w := NewWorker(1, queue, registry)

	// Assert
	assert.Equal(t, 1, w.id)
	assert.Equal(t, queue, w.queue)
	assert.Equal(t, registry, w.taskHandlerRegistry)
}

func TestWorker_Start(t *testing.T) {
	t.Run("Processes a task successfully", func(t *testing.T) {
		// Arrange
		queue := make(chan task.Task)
		registry := new(mockTaskHandlerRegistry)

		taskToProcess := task.Task{
			ID:      "task-1",
			Type:    "test-task",
			Payload: "test-payload",
		}

		var receivedPayload string

		handler := task.Handler(func(payload any) error {
			receivedPayload, _ = payload.(string)
			return nil
		})

		registry.On("GetHandler", taskToProcess.Type).Return(handler, true).Once()

		w := NewWorker(1, queue, registry)

		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup

		wg.Add(1)

		// Act
		w.Start(ctx, &wg)
		queue <- taskToProcess

		cancel()

		// Assert
		wg.Wait()

		registry.AssertExpectations(t)

		assert.Equal(t, taskToProcess.Payload, receivedPayload)
	})

	t.Run("Logs warning and skips task when no handler is registered", func(t *testing.T) {
		// Arrange
		var logOutput bytes.Buffer

		logrus.SetOutput(&logOutput) // Capture log output
		logrus.SetLevel(logrus.WarnLevel)

		queue := make(chan task.Task)
		registry := new(mockTaskHandlerRegistry)

		taskToProcess := task.Task{
			ID:      "task-1",
			Type:    "test-task",
			Payload: "test-payload",
		}

		registry.On("GetHandler", taskToProcess.Type).Return(nil, false).Once()

		w := NewWorker(1, queue, registry)

		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup

		wg.Add(1)

		// Act
		w.Start(ctx, &wg)
		queue <- taskToProcess

		cancel()

		// Assert
		wg.Wait()

		registry.AssertExpectations(t)

		require.Contains(t, logOutput.String(), "No handler registered for task type")
		require.Contains(t, logOutput.String(), "worker_id=1")
		require.Contains(t, logOutput.String(), "task_type=test-task")
	})

	t.Run("Logs error when handler reports error", func(t *testing.T) {
		// Arrange
		var logOutput bytes.Buffer

		logrus.SetOutput(&logOutput) // Capture log output
		logrus.SetLevel(logrus.ErrorLevel)

		queue := make(chan task.Task, 1)
		registry := new(mockTaskHandlerRegistry)

		taskToProcess := task.Task{
			ID:      "task-1",
			Type:    "test-task",
			Payload: "test-payload",
		}

		errorToReturn := errors.New("error")

		var receivedPayload string

		handler := task.Handler(func(payload any) error {
			receivedPayload, _ = payload.(string)
			return errorToReturn
		})

		registry.On("GetHandler", taskToProcess.Type).Return(handler, true).Once()

		w := NewWorker(1, queue, registry)

		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup

		wg.Add(1)

		// Act
		w.Start(ctx, &wg)

		queue <- taskToProcess

		cancel()

		// Assert
		wg.Wait()

		registry.AssertExpectations(t)

		assert.Equal(t, taskToProcess.Payload, receivedPayload)

		require.Contains(t, logOutput.String(), "Failed to process task")
		require.Contains(t, logOutput.String(), "worker_id=1")
		require.Contains(t, logOutput.String(), "task_id=task-1")
		require.Contains(t, logOutput.String(), fmt.Sprintf("error=%s", errorToReturn))
	})
}
