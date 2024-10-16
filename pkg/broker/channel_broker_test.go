package broker

import (
	"context"
	"testing"
	"time"

	"github.com/jamesTait-jt/goflow/task"
	"github.com/stretchr/testify/assert"
)

func Test_NewChannelBroker(t *testing.T) {
	t.Run("Creates a channel broker with the correct buffer size", func(t *testing.T) {
		// Arrange
		bufferSize := 10

		// Act
		b := NewChannelBroker(bufferSize)

		// Assert
		assert.Equal(t, bufferSize, cap(b.taskQueue))
	})
}

func Test_ChannelBroker_Submit(t *testing.T) {
	t.Run("Places the task on the task queue", func(t *testing.T) {
		// Arrange
		tsk := task.Task{
			ID: "randomID",
		}

		b := ChannelBroker{
			taskQueue: make(chan task.Task, 1),
		}

		// Act
		b.Submit(context.Background(), tsk)

		// Assert
		assert.Equal(t, tsk, <-b.taskQueue)
	})

	t.Run("Exits when the context is cancelled", func(t *testing.T) {
		// Arrange
		b := ChannelBroker{
			taskQueue: make(chan task.Task),
		}

		// Use a timeout context to ensure the test doesn't run forever in case of failure
		testCtx, testCancel := context.WithTimeout(context.Background(), time.Second)
		defer testCancel()

		ctx, cancel := context.WithCancel(context.Background())

		// Act
		cancel()

		done := make(chan struct{})
		go func() {
			b.Submit(ctx, task.Task{})
			close(done)
		}()

		// Assert
		select {
		case <-done:
			// Success, the goroutine returned as expected
		case <-testCtx.Done():
			// The test context timed out, meaning the Submit method didn't exit as expected
			t.Fatal("Submit did not return after context was cancelled")
		}
	})

	t.Run("Exits when the context is cancelled after submitting blocked task", func(t *testing.T) {
		// Arrange
		b := ChannelBroker{
			taskQueue: make(chan task.Task),
		}

		// Use a timeout context to ensure the test doesn't run forever in case of failure
		testCtx, testCancel := context.WithTimeout(context.Background(), time.Second)
		defer testCancel()

		ctx, cancel := context.WithCancel(context.Background())

		// Act
		done := make(chan struct{})
		go func() {
			b.Submit(ctx, task.Task{})
			close(done)
		}()

		cancel()

		// Assert
		select {
		case <-done:
			// Success, the goroutine returned as expected
		case <-testCtx.Done():
			// The test context timed out, meaning the Submit method didn't exit as expected
			t.Fatal("Submit did not return after context was cancelled")
		}
	})
}

func Test_ChannelBroker_Dequeue(t *testing.T) {
	t.Run("Returns the task queue", func(t *testing.T) {
		// Arrange
		tsk := task.Task{
			ID: "randomID",
		}

		b := ChannelBroker{
			taskQueue: make(chan task.Task, 1),
		}

		b.taskQueue <- tsk

		// Act
		taskQueue := b.Dequeue(context.Background())

		// Assert
		assert.Equal(t, tsk, <-taskQueue)
	})
}
