package worker

import (
	"context"
	"testing"

	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockTaskHandlerRegistry is a mock for the task registry interface
type mockTaskHandlerRegistry struct {
	mock.Mock
}

func (m *mockTaskHandlerRegistry) GetHandler(taskType string) (task.TaskHandler, bool) {
	args := m.Called(taskType)
	return args.Get(0).(task.TaskHandler), args.Bool(1)
}

func TestNewWorker_Creates_a_new_worker_with_variables_initialised(t *testing.T) {
	// Arrange
	queue := make(<-chan task.Task, 1)
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
		queue := make(chan task.Task, 1)
		registry := new(mockTaskHandlerRegistry)

		taskToProcess := task.Task{
			ID:      "task-1",
			Type:    "test-task",
			Payload: "test-payload",
		}

		receivedPayloads := make(chan any, 1)
		handler := task.TaskHandler(func(payload any) error {
			receivedPayloads <- payload
			return nil
		})

		registry.On("GetHandler", taskToProcess.Type).Return(handler, true).Once()

		w := NewWorker(1, queue, registry)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Act
		w.Start(ctx)
		queue <- taskToProcess

		// Assert
		receivedPayload := <-receivedPayloads
		assert.Equal(t, taskToProcess.Payload, receivedPayload)
		registry.AssertExpectations(t)
	})
}
