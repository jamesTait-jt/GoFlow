package task

import (
	"context"

	"github.com/google/uuid"
)

// Handler processes task payloads
type Handler func(payload any) Result

// Type represents a generic task structure
type Task struct {
	ID      string
	Type    string
	Payload any
}

type Result struct {
	TaskID  string
	Payload any
	Error   error
}

func New(taskType string, payload any) Task {
	id := uuid.New()
	t := Task{
		ID:      id.String(),
		Type:    taskType,
		Payload: payload,
	}

	return t
}

type TaskOrResult interface {
	Task | Result
}

type Submitter[T TaskOrResult] interface {
	Submit(ctx context.Context, t T)
}

type Dequeuer[T TaskOrResult] interface {
	Dequeue(ctx context.Context) <-chan T
}
