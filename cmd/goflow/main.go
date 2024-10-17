package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) Greet(ctx context.Context, req *GreetRequest) (*GreetResponse, error) {
	return &GreetResponse{Message: "Hello " + req.Name}, nil
}

// Request structure for the Greet method
type GreetRequest struct {
	Name string
}

// Response structure for the Greet method
type GreetResponse struct {
	Message string
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	registerGreetServiceServer(grpcServer, &Server{})

	fmt.Println("gRPC server running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func registerGreetServiceServer(s *grpc.Server, srv *Server) {
	// Normally, this would register the service methods
	// For now, just a placeholder
}
