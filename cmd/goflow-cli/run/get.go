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

func Get(taskID string) error {
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

	r, err := goFlowClient.GetResult(ctx, &pb.GetResultRequest{TaskID: taskID})
	if err != nil {
		log.Printf("could not get result: %v", err)
		return nil
	}

	log.Printf("Result: %s", r.GetResult())

	return nil
}
