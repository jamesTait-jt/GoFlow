package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesTait-jt/goflow/pkg/task"
	"github.com/redis/go-redis/v9"
)

type RedisBroker[T task.TaskOrResult] struct {
	client        *redis.Client
	redisQueueKey string
}

func NewRedisBroker[T task.TaskOrResult](
	client *redis.Client, key string,
) *RedisBroker[T] {
	return &RedisBroker[T]{
		client:        client,
		redisQueueKey: key,
	}
}

func (rb *RedisBroker[T]) Submit(ctx context.Context, submission T) error {
	serialised, err := task.Serialize(submission)
	if err != nil {
		return err
	}

	_, err = rb.client.LPush(ctx, rb.redisQueueKey, serialised).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rb *RedisBroker[T]) Dequeue(ctx context.Context) <-chan T {
	ch := make(chan T, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				redisResult, err := rb.client.BRPop(ctx, time.Second, rb.redisQueueKey).Result()
				if err != nil {
					if err == redis.Nil {
						// BRPop timed out
						continue
					}
					fmt.Println("BRPop error: " + err.Error())
					continue
				}

				result, err := task.Deserialize[T]([]byte(redisResult[1]))
				if err != nil {
					fmt.Println("Failed to deserialize task:", err)

					continue
				}

				ch <- result
			}
		}
	}()

	return ch
}
