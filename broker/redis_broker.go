package broker

import (
	"context"

	"github.com/jamesTait-jt/goflow/task"
	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client        redis.Client
	redisQueueKey string
}

func NewRedisBroker(client redis.Client) *RedisBroker {
	return &RedisBroker{
		client:        client,
		redisQueueKey: "queue",
	}
}

func (rb *RedisBroker) Submit(ctx context.Context, t task.Task) {
	rb.client.LPush(ctx, rb.redisQueueKey, t)
}
