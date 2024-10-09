package broker

// import (
// 	"github.com/jamesTait-jt/GoFlow/internal/task"
// 	"github.com/redis/go-redis/v9"
// )

// type RedisBroker struct {
// 	client      *redis.Client
// 	taskChannel chan<- task.Task
// }

// func NewRedisBroker(redisAddr string, taskChannel chan<- task.Task) *RedisBroker {
// 	client := redis.NewClient(&redis.Options{Addr: redisAddr})

// 	broker := &RedisBroker{
// 		client:      client,
// 		taskChannel: taskChannel,
// 	}
// }

// func (rd *RedisBroker) Submit(t task.Task) {

// }
