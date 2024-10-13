package broker

import (
	"testing"

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
		b.Submit(tsk)

		// Assert
		assert.Equal(t, tsk, <-b.taskQueue)
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
		taskQueue := b.Dequeue()

		// Assert
		assert.Equal(t, tsk, <-taskQueue)
	})
}
