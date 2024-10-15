package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesTait-jt/goflow/task"
	"github.com/redis/go-redis/v9"
)

type RedisBroker[S task.TaskOrResult, D task.TaskOrResult] struct {
	client        *redis.Client
	redisQueueKey string
}

func NewRedisBroker[S task.TaskOrResult, D task.TaskOrResult](
	client *redis.Client, key string,
) *RedisBroker[S, D] {
	return &RedisBroker[S, D]{
		client:        client,
		redisQueueKey: key,
	}
}

func (rb *RedisBroker[S, D]) Submit(ctx context.Context, submission S) error {
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

func (rb *RedisBroker[S, D]) Dequeue(ctx context.Context) <-chan D {
	ch := make(chan D, 1)

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

				result, err := task.Deserialize[D]([]byte(redisResult[1]))
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
