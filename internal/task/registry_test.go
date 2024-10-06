package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskHandlerRegistry_RegisterTaskHandler(t *testing.T) {
	t.Run("registers the new task handler when the task type was not registered", func(t *testing.T) {
		// Arrange
		registry := NewHandlerRegistry()

		invoked := false
		mockHandler := func(_ any) error {
			invoked = true
			return nil
		}

		taskType := "newTaskType"

		// Act
		registry.RegisterTaskHandler(taskType, mockHandler)

		// Assert
		addedHandler, ok := registry.handlers[taskType]

		assert.True(t, ok)

		_ = addedHandler(nil)

		assert.True(t, invoked)
	})

	t.Run("overwrites the existing task handler for the same task type", func(t *testing.T) {
		// Arrange
		registry := NewHandlerRegistry()
		taskType := "newTaskType"

		invokedOriginal := false
		originalHandler := func(_ any) error {
			invokedOriginal = true
			return nil
		}

		registry.handlers[taskType] = originalHandler

		invokedNew := false
		newHandler := func(_ any) error {
			invokedNew = true
			return nil
		}

		// Act
		registry.RegisterTaskHandler(taskType, newHandler)

		// Assert
		addedHandler, ok := registry.handlers[taskType]

		assert.True(t, ok)

		_ = addedHandler(nil)

		assert.True(t, invokedNew)
		assert.False(t, invokedOriginal)
	})
}

func TestTaskHandlerRegistry_GetHandler(t *testing.T) {
	t.Run("retrieves nil when the task type is unregistered", func(t *testing.T) {
		// Arrange
		registry := NewHandlerRegistry()

		// Act
		handler, ok := registry.GetHandler("unregisteredTaskType")

		// Assert
		assert.False(t, ok)
		assert.Nil(t, handler)
	})

	t.Run("retrieves the correct handler when the task type is registered", func(t *testing.T) {
		// Arrange
		registry := NewHandlerRegistry()
		taskType := "newTaskType"

		invoked := false
		mockHandler := func(_ any) error {
			invoked = true
			return nil
		}

		registry.handlers[taskType] = mockHandler

		// Act
		handler, ok := registry.GetHandler(taskType)

		// Assert
		assert.True(t, ok)

		_ = handler(nil)

		assert.True(t, invoked)
	})
}
