package dispatcher

import (
	"github.com/jamesTait-jt/GoFlow/internal/task"
	"github.com/redis/go-redis/v9"
)

type RedisDispatcher struct {
	client      *redis.Client
	taskChannel chan<- task.Task
}

func NewRedisDispatcher(redisAddr str, taskChannel chan<- task.Task) *RedisDispatcher {
	client := redis.NewClient(&redis.Options{Addr: redisAddr})
}

func (rd *RedisDispatcher) Dispatch(t task.Task) {

}
