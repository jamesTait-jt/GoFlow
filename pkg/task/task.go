package task

import (
	"bytes"
	"context"
	"encoding/gob"

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
	Submit(ctx context.Context, t T) error
}

type Dequeuer[T TaskOrResult] interface {
	Dequeue(ctx context.Context) <-chan T
}

func Serialize[T any](t T) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(t)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize[T any](data []byte) (T, error) {
	var t T

	buf := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(&t)

	if err != nil {
		return t, err
	}

	return t, nil
}
