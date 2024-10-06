package task

// Type represents a generic task structure
type Task struct {
	ID      string
	Type    string
	Payload any
}

// Handler processes task payloads
type Handler func(payload any) error
