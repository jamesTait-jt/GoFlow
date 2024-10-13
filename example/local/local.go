package main

import (
	"fmt"
	"time"

	goflow "github.com/jamesTait-jt/goflow"
	"github.com/jamesTait-jt/goflow/broker"
	"github.com/jamesTait-jt/goflow/store"
	"github.com/jamesTait-jt/goflow/task"
	"github.com/jamesTait-jt/goflow/worker"
	"github.com/jamesTait-jt/goflow/workerpool"
)

func main() {
	var workers []workerpool.Worker
	for i := 0; i < 5; i++ {
		workers = append(workers, worker.NewLocalWorker())
	}

	taskHandlerStore := store.NewInMemoryKVStore[string, task.Handler]()
	resultsStore := store.NewInMemoryKVStore[string, task.Result]()
	channelBroker := broker.NewChannelBroker(5)

	gf := goflow.New(
		workers,
		goflow.WithTaskHandlerStore(taskHandlerStore),
		goflow.WithResultsStore(resultsStore),
		goflow.WithTaskBroker(channelBroker),
	)

	// Example task handler
	taskHandler := func(payload any) task.Result {
		// Simulate some processing
		return task.Result{Payload: fmt.Sprintf("Processed: %v", payload)}
	}

	// Register the handler
	taskType := "exampleTask"
	gf.RegisterHandler(taskType, taskHandler)

	gf.Start()

	// Push a task to the goflow
	taskIDs := []string{}

	for i := 0; i < 10; i++ {
		taskID, err := gf.Push(taskType, "My example payload")
		if err != nil {
			fmt.Printf("Error pushing task: %v\n", err)
			return
		}

		fmt.Printf("Task submitted with ID: %s\n", taskID)
		taskIDs = append(taskIDs, taskID)
	}

	time.Sleep(time.Second * 1)

	for i := 0; i < len(taskIDs); i++ {
		result, _ := gf.GetResult(taskIDs[i])
		fmt.Println(result)
	}
}
