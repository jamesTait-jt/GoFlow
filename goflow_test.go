package goflow

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jamesTait-jt/GoFlow/task"
	"github.com/jamesTait-jt/GoFlow/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWorker struct {
	mock.Mock
}

func (m *mockWorker) Start(ctx context.Context, wg *sync.WaitGroup, taskSource worker.TaskSource) {
	m.Called(ctx, wg, taskSource)
}

type mockTaskBroker struct {
	mock.Mock
}

func (m *mockTaskBroker) Submit(t task.Task) {
	m.Called(t)
}

func (m *mockTaskBroker) Dequeue() <-chan task.Task {
	args := m.Called()
	return args.Get(0).(<-chan task.Task)
}

type mockKVStore[K comparable, V any] struct {
	mock.Mock
}

func (m *mockKVStore[K, V]) Put(key K, value V) {
	m.Called(key, value)
}

func (m *mockKVStore[K, V]) Get(key K) (V, bool) {
	args := m.Called(key)
	return args.Get(0).(V), args.Bool(1)
}

func Test_NewGoFlow(t *testing.T) {
	t.Run("Initialises GoFlow properly", func(t *testing.T) {
		// Arrange
		mockTaskBroker := new(mockTaskBroker)
		mockHandlers := new(mockKVStore[string, task.Handler])
		mockResults := new(mockKVStore[string, task.Result])

		// Act
		gf := NewGoFlow(0, nil, mockHandlers, mockResults, mockTaskBroker)

		// Assert
		assert.NotNil(t, gf)
		assert.NotNil(t, gf.ctx)
		assert.NotNil(t, gf.cancel)
		assert.NotNil(t, gf.workers)

		assert.Equal(t, mockTaskBroker, gf.taskBroker)
		assert.Equal(t, mockHandlers, gf.taskHandlers)
		assert.Equal(t, mockResults, gf.results)
	})

	t.Run("Initialises the worker pool with correct number of workers", func(t *testing.T) {
		// Arrange
		ids := []int{}
		mockWorkerFactory := func(id int) worker.Worker {
			ids = append(ids, id)
			return new(mockWorker)
		}

		// Act
		NewGoFlow(3, mockWorkerFactory, nil, nil, nil)

		// Assert
		assert.Equal(t, []int{0, 1, 2}, ids)
	})
}

func Test_GoFlow_RegisterHandler(t *testing.T) {
	t.Run("Puts the handler in the handler store", func(t *testing.T) {
		// Arrange
		mockHandlers := new(mockKVStore[string, task.Handler])
		handler := func(_ any) task.Result {
			return task.Result{}
		}
		gf := GoFlow{
			taskHandlers: mockHandlers,
		}
		taskType := "exampleTask"

		mockHandlers.On("Put", taskType, mock.AnythingOfType("task.Handler")).Once()

		// Act
		gf.RegisterHandler(taskType, handler)

		// Assert
		mockHandlers.AssertExpectations(t)
	})
}

func Test_GoFlow_Push(t *testing.T) {
	t.Run("Returns an error if the handler is not registered", func(t *testing.T) {
		// Arrange
		mockHandlers := new(mockKVStore[string, task.Handler])
		gf := GoFlow{
			taskHandlers: mockHandlers,
		}
		taskType := "exampleTask"

		var mockHandler task.Handler = func(_ any) task.Result {
			return task.Result{}
		}

		mockHandlers.On("Get", taskType).Once().Return(mockHandler, false)

		// Act
		taskID, err := gf.Push(taskType, "payload")

		// Assert
		assert.Equal(t, "", taskID)
		assert.EqualError(t, err, fmt.Sprintf("no handler defined for taskType: %s", taskType))
	})
}
