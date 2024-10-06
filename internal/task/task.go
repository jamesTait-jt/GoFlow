package task

// Type represents a generic task structure
type Task struct {
	ID      string
	Type    string
	Payload any
}

// TaskHandler processes task payloads
type TaskHandler func(payload any) error
