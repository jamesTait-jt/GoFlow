package main

import (
	"fmt"
	"time"

	"github.com/jamesTait-jt/goflow/pkg/task"
	"golang.org/x/exp/rand"
)

func handle(payload any) task.Result {
	rand.Seed(uint64(time.Now().UnixNano()))
	n := rand.Intn(1000) // n will be between 0 and 10
	fmt.Printf("Sleeping %d milliseconds...\n", n)
	time.Sleep(time.Millisecond * time.Duration(n))
	fmt.Println("Done")
	fmt.Println("Handling task with payload: ", payload)

	return task.Result{}
}

func NewHandler() task.Handler {
	return handle
}
