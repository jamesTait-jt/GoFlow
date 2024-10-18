package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	pb "github.com/jamesTait-jt/goflow/cmd/goflow/goflow"
	goflow "github.com/jamesTait-jt/goflow/pkg"
	"github.com/jamesTait-jt/goflow/pkg/broker"
	"github.com/jamesTait-jt/goflow/pkg/store"
	"github.com/jamesTait-jt/goflow/pkg/task"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var (
	redisPort  = "6379"
	serverPort = "50051"
)

type server struct {
	pb.GoFlowServer
	gf *goflow.GoFlow
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) PushTask(_ context.Context, in *pb.PushTaskRequest) (*pb.PushTaskReply, error) {
	log.Printf("Received push task: [%s] [%s]", in.GetTaskType(), in.GetPayload())

	id, err := s.gf.Push(in.GetTaskType(), in.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to push task: %v", err)
	}

	return &pb.PushTaskReply{Id: id}, nil
}

func (s *server) GetResult(_ context.Context, in *pb.GetResultRequest) (*pb.GetResultReply, error) {
	log.Printf("Received get result: [%s]", in.GetTaskID())

	result, ok := s.gf.GetResult(in.GetTaskID())
	if !ok {
		return nil, fmt.Errorf("task not complete or didnt exist")
	}

	parsedResult, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", result)
	}

	return &pb.GetResultReply{Result: string(parsedResult)}, nil
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("goflow-redis-server:%s", redisPort),
	})
	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	log.Printf("redis connection successful: %s", pong)

	taskSubmitter := broker.NewRedisBroker[task.Task](redisClient, "tasks")
	resultsGetter := broker.NewRedisBroker[task.Result](redisClient, "results")
	resultsStore := store.NewInMemoryKVStore[string, task.Result]()

	gf := goflow.New(
		goflow.WithTaskBroker(taskSubmitter),
		goflow.WithResultBroker(resultsGetter),
		goflow.WithResultsStore(resultsStore),
	)

	gf.Start()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoFlowServer(grpcServer, &server{gf: gf})

	log.Printf("server listening at %v", lis.Addr())
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
