package run

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jamesTait-jt/goflow/cmd/goflow-cli/config"
	pb "github.com/jamesTait-jt/goflow/cmd/goflow/goflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Push(taskType, payload string) error {
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", config.GoFlowHostPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	goFlowClient := pb.NewGoFlowClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := goFlowClient.PushTask(ctx, &pb.PushTaskRequest{TaskType: taskType, Payload: payload})
	if err != nil {
		log.Fatalf("could not push: %v", err)
	}
	log.Printf("Task ID: %s", r.GetId())

	return nil
}
