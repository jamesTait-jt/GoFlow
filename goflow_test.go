package goflow

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jamesTait-jt/goflow/broker"
	"github.com/jamesTait-jt/goflow/store"
	"github.com/jamesTait-jt/goflow/task"
	"github.com/jamesTait-jt/goflow/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWorkerPool struct {
	mock.Mock
}

func (m *mockWorkerPool) Start(ctx context.Context, taskSource workerpool.TaskSource) {
	m.Called(ctx, taskSource)
}

func (m *mockWorkerPool) AwaitShutdown() {
	m.Called()
}

type mockTaskBroker struct {
	mock.Mock
}

func (m *mockTaskBroker) Submit(ctx context.Context, t task.Task) {
	m.Called(ctx, t)
}

func (m *mockTaskBroker) Dequeue(ctx context.Context) <-chan task.Task {
	args := m.Called(ctx)
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

func Test_New(t *testing.T) {
	t.Run("Initialises goflow with default options", func(t *testing.T) {
		// Arrange
		// Act
		gf := New([]workerpool.Worker{})

		// Assert
		assert.NotNil(t, gf)
		assert.NotNil(t, gf.ctx)
		assert.NotNil(t, gf.cancel)
		assert.NotNil(t, gf.workers)

		assert.IsType(t, &broker.ChannelBroker{}, gf.taskBroker)
		assert.IsType(t, &store.InMemoryKVStore[string, task.Handler]{}, gf.taskHandlers)
		assert.IsType(t, &store.InMemoryKVStore[string, task.Result]{}, gf.results)
	})
	t.Run("Initialises goflow with custom options", func(t *testing.T) {
		// Arrange
		mockTaskBroker := new(mockTaskBroker)
		mockHandlers := new(mockKVStore[string, task.Handler])
		mockResults := new(mockKVStore[string, task.Result])

		// Act
		gf := New(
			[]workerpool.Worker{},
			WithTaskBroker(mockTaskBroker),
			WithTaskHandlerStore(mockHandlers),
			WithResultsStore(mockResults),
		)

		// Assert
		assert.NotNil(t, gf)
		assert.NotNil(t, gf.ctx)
		assert.NotNil(t, gf.cancel)
		assert.NotNil(t, gf.workers)

		assert.Equal(t, mockTaskBroker, gf.taskBroker)
		assert.Equal(t, mockHandlers, gf.taskHandlers)
		assert.Equal(t, mockResults, gf.results)
	})
}

func Test_goflow_RegisterHandler(t *testing.T) {
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

func Test_goflow_Push(t *testing.T) {
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

	t.Run("Submits the task to the broker and handles the result", func(t *testing.T) {
		// Arrange
		mockHandlers := new(mockKVStore[string, task.Handler])
		mockBroker := new(mockTaskBroker)
		mockResults := new(mockKVStore[string, task.Result])

		ctx := context.Background()

		gf := GoFlow{
			ctx:          ctx,
			taskHandlers: mockHandlers,
			taskBroker:   mockBroker,
			results:      mockResults,
		}

		var mockHandler task.Handler = func(_ any) task.Result {
			return task.Result{}
		}

		taskType := "exampleTask"
		payload := "examplePayload"
		responseFromWorker := task.Result{Payload: "Successful work!!"}

		mockHandlers.On("Get", mock.Anything).Once().Return(mockHandler, true)

		var submittedTask task.Task

		mockBroker.On("Submit", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {
			submittedTask, _ = args.Get(1).(task.Task)

			// Put something on the created resultChannel to simulate a response from the worker
			submittedTask.ResultCh <- responseFromWorker
		})

		var wg sync.WaitGroup

		wg.Add(1)

		mockResults.On("Put", mock.Anything, mock.Anything).Once().Run(func(_ mock.Arguments) {
			// Signal that the result was persisted, allowing the test to continue
			wg.Done()
		})

		// Act
		taskID, err := gf.Push(taskType, payload)

		wg.Wait()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, submittedTask.ID, taskID)

		assert.Equal(t, taskType, submittedTask.Type)
		assert.Equal(t, payload, submittedTask.Payload)

		mockHandlers.AssertCalled(t, "Get", taskType)
		mockResults.AssertCalled(t, "Put", submittedTask.ID, responseFromWorker)
		mockBroker.AssertCalled(t, "Submit", ctx, mock.AnythingOfType("task.Task"))
	})
}

func Test_goflow_GetResult(t *testing.T) {
	t.Run("Returns the result of given taskID if it exists", func(t *testing.T) {
		// Arrange
		mockResults := new(mockKVStore[string, task.Result])

		gf := GoFlow{
			results: mockResults,
		}

		taskID := "taskID"

		expectedResult := task.Result{Payload: "result"}

		mockResults.On("Get", mock.Anything).Once().Return(expectedResult, true)

		// Act
		result, ok := gf.GetResult(taskID)

		// Assert
		assert.Equal(t, expectedResult, result)
		assert.True(t, ok)
	})
	t.Run("Returns false if given taskID doesn't exist", func(t *testing.T) {
		// Arrange
		mockResults := new(mockKVStore[string, task.Result])

		gf := GoFlow{
			results: mockResults,
		}

		taskID := "taskID"

		expectedResult := task.Result{}

		mockResults.On("Get", mock.Anything).Once().Return(expectedResult, false)

		// Act
		result, ok := gf.GetResult(taskID)

		// Assert
		assert.Equal(t, expectedResult, result)
		assert.False(t, ok)
	})
}

func Test_goflow_Stop(t *testing.T) {
	t.Run("Calls cancel and waits for all workers to shut down", func(t *testing.T) {
		// Arrange
		wasCancelCalled := false
		mockCancel := func() {
			wasCancelCalled = true
		}

		mockWorkerPool := &mockWorkerPool{}
		mockWorkerPool.On("AwaitShutdown").Once()

		gf := GoFlow{
			cancel:  mockCancel,
			workers: mockWorkerPool,
		}

		// Act
		gf.Stop()

		// Assert
		assert.True(t, wasCancelCalled)
		mockWorkerPool.AssertExpectations(t)
	})
}
