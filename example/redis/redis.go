package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jamesTait-jt/goflow/pkg/broker"
	"github.com/jamesTait-jt/goflow/pkg/task"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	redisBroker := broker.NewRedisBroker[task.Task](redisClient, "tasks")

	ctx, _ := context.WithCancel(context.Background())

	for i := 0; i < 100; i++ {
		t := task.Task{ID: strconv.Itoa(i), Type: "testplugin", Payload: fmt.Sprintf("im a random sleeper")}
		redisBroker.Submit(ctx, t)
	}
}
