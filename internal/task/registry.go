package task

// HandlerRegistry is a map for storing task handlers
type HandlerRegistry struct {
	handlers map[string]Handler
}

// NewHandlerRegistry creates a new TaskRegistry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]Handler),
	}
}

// RegisterTaskHandler registers a handler for a given task type, overwriting any existing handlers
func (thr *HandlerRegistry) RegisterHandler(taskType string, handler Handler) {
	thr.handlers[taskType] = handler
}

// GetHandler retrieves the correct handler based on the task type
func (thr *HandlerRegistry) GetHandler(taskType string) (Handler, bool) {
	handler, ok := thr.handlers[taskType]

	return handler, ok
}
