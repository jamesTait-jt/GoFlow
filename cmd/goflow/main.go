package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/jamesTait-jt/goflow/cmd/goflow/goflow"
	"google.golang.org/grpc"
)

type server struct {
	pb.GoFlowServer
}

func (s *server) PushTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	return &pb.TaskResponse{Message: "Hello " + req.TaskType + " " + req.Payload}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoFlowServer(grpcServer, &server{})

	fmt.Println("gRPC server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
