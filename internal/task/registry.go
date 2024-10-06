package task

// TaskHandlerRegistry is a map for storing task handlers
type TaskHandlerRegistry struct {
	handlers map[string]TaskHandler
}

// NewTaskHandlerRegistry creates a new TaskRegistry
func NewTaskHandlerRegistry() *TaskHandlerRegistry {
	return &TaskHandlerRegistry{
		handlers: make(map[string]TaskHandler),
	}
}

// RegisterTaskHandler registers a handler for a given task type, overwriting any existing handlers
func (thr *TaskHandlerRegistry) RegisterTaskHandler(taskType string, handler TaskHandler) {
	thr.handlers[taskType] = handler
}

// GetHandler retrieves the correct handler based on the task type
func (thr *TaskHandlerRegistry) GetHandler(taskType string) (TaskHandler, bool) {
	handler, ok := thr.handlers[taskType]

	return handler, ok
}
