package main

import (
	"fmt"

	"github.com/jamesTait-jt/goflow/pkg/task"
)

func handle(payload any) task.Result {
	fmt.Println("Handling task with payload: ", payload)

	return task.Result{}
}

func NewHandler() task.Handler {
	return handle
}
