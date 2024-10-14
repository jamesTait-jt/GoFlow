package task

import (
	"github.com/google/uuid"
)

// Handler processes task payloads
type Handler func(payload any) Result

// Type represents a generic task structure
type Task struct {
	ID      string
	Type    string
	Payload any
	Handler Handler
}

type Result struct {
	TaskID  string
	Payload any
	Error   error
}

func New(taskType string, payload any, handler Handler) Task {
	id := uuid.New()
	t := Task{
		ID:      id.String(),
		Type:    taskType,
		Payload: payload,
		Handler: handler,
	}

	return t
}
