package main

import (
	"fmt"
	"time"

	goflow "github.com/jamesTait-jt/goflow/pkg"
	"github.com/jamesTait-jt/goflow/pkg/broker"
	"github.com/jamesTait-jt/goflow/pkg/store"
	"github.com/jamesTait-jt/goflow/pkg/task"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	taskSubmitter := broker.NewRedisBroker[task.Task](redisClient, "tasks")
	resultsGetter := broker.NewRedisBroker[task.Result](redisClient, "results")
	resultsStore := store.NewInMemoryKVStore[string, task.Result]()

	gf := goflow.New(
		goflow.WithTaskBroker(taskSubmitter),
		goflow.WithResultBroker(resultsGetter),
		goflow.WithResultsStore(resultsStore),
	)

	gf.Start()

	results := make(chan task.Result)
	for i := 0; i < 100; i++ {
		id, _ := gf.Push("testplugin", "Im a random sleeper")
		go func() {
			for {
				result, ok := gf.GetResult(id)
				if !ok {
					time.Sleep(time.Second)
					continue
				}

				results <- result
				if i == 99 {
					close(results)
				}

				return
			}
		}()
	}

	i := 0
	for r := range results {
		i++
		fmt.Println(i, r)
	}
}
