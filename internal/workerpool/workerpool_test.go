package workerpool

import (
	"context"
	"sync"
	"testing"

	"github.com/jamesTait-jt/GoFlow/task"
	"github.com/jamesTait-jt/GoFlow/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWorker struct {
	mock.Mock
}

func (m *mockWorker) Start(ctx context.Context, wg *sync.WaitGroup, taskSource workerpool.TaskSource) {
	m.Called(ctx, wg, taskSource)
}

type mockTaskSource struct {
	mock.Mock
}

func (m *mockTaskSource) Dequeue() <-chan task.Task {
	args := m.Called()
	return args.Get(0).(<-chan task.Task)
}

func TestNewWorkerPool(t *testing.T) {
	t.Run("Creates a new worker pool with variables initialised", func(t *testing.T) {
		// Arrange
		numWorkers := 5

		ids := []int{}
		workerFactory := func(id int) workerpool.Worker {
			ids = append(ids, id)
			return &mockWorker{}
		}

		// Act
		wp := NewWorkerPool(numWorkers, workerFactory)

		// Assert
		assert.Equal(t, numWorkers, len(wp.workers))
		assert.NotNil(t, wp.wg)
		assert.Equal(t, []int{0, 1, 2, 3, 4}, ids)
	})
}

func TestPool_Start(t *testing.T) {
	t.Run("Starts all workers in the pool", func(t *testing.T) {
		// Arrange
		numWorkers := 5
		ctx := context.Background()
		wg := &sync.WaitGroup{}
		taskSource := &mockTaskSource{}

		mockWorkers := make(map[int]*mockWorker)

		for i := 0; i < numWorkers; i++ {
			mockWorker := new(mockWorker)
			mockWorker.On("Start", ctx, wg, taskSource).Once()
			mockWorkers[i] = mockWorker
		}

		pool := &Pool{
			workers: make(map[int]workerpool.Worker),
			wg:      wg,
		}

		for i := 0; i < numWorkers; i++ {
			pool.workers[i] = mockWorkers[i]
		}

		// Act
		pool.Start(ctx, taskSource)

		// Assert
		for i := 0; i < numWorkers; i++ {
			mockWorkers[i].AssertCalled(t, "Start", ctx, wg, taskSource)
		}
	})
}
