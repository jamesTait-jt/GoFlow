package worker

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWorker struct {
	mock.Mock
}

func (m *mockWorker) Start(ctx context.Context, wg *sync.WaitGroup) {
	m.Called(ctx, wg)
}

func TestNewWorkerPool(t *testing.T) {
	t.Run("Creates a new worker pool with variables initialised", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		numWorkers := 5

		// Act
		wp := NewWorkerPool(ctx, numWorkers, nil)

		// Assert
		assert.Equal(t, numWorkers, len(wp.workers))
		assert.Equal(t, ctx, wp.ctx)
		assert.NotNil(t, wp.wg)
	})
}

func TestPool_Start(t *testing.T) {
	t.Run("Starts all workers in the pool", func(t *testing.T) {
		// Arrange
		numWorkers := 5
		ctx := context.Background()
		wg := &sync.WaitGroup{}

		mockWorkers := make(map[int]*mockWorker)

		for i := 0; i < numWorkers; i++ {
			mockWorker := new(mockWorker)
			mockWorker.On("Start", ctx, wg).Once()
			mockWorkers[i] = mockWorker
		}

		pool := &Pool{
			workers: make(map[int]worker),
			ctx:     ctx,
			wg:      wg,
		}

		for i := 0; i < numWorkers; i++ {
			pool.workers[i] = mockWorkers[i]
		}

		// Act
		pool.Start()

		// Assert
		for i := 0; i < numWorkers; i++ {
			mockWorkers[i].AssertCalled(t, "Start", ctx, wg)
		}
	})
}
