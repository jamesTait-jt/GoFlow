package main

import (
	"context"
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

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("redis:%s", redisPort),
	})
	log.Print("connected to redis")

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
