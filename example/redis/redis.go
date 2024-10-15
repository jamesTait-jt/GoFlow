package main

import (
	"context"

	"github.com/jamesTait-jt/goflow/broker"
	"github.com/jamesTait-jt/goflow/task"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	redisBroker := broker.NewRedisBroker[task.Task, task.Task](redisClient, "tasks")

	ctx, _ := context.WithCancel(context.Background())

	t := task.Task{Type: "testplugin", Payload: "payload"}
	redisBroker.Submit(ctx, t)
}
