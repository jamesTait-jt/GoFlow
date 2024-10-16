package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var dockerNetworkName = "goflow-network"
var redisContainerName = "redis-server"

func main() {
	// Define Cobra root command
	rootCmd := &cobra.Command{
		Use:   "goflow",
		Short: "Goflow CLI tool to deploy workerpool and plugins using Docker",
	}

	// Define deploy command
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy workerpool with Redis broker and compiled plugins",
		Run: func(cmd *cobra.Command, args []string) {
			err := deploy()
			if err != nil {
				log.Fatalf("Error during deployment: %v", err)
			}
		},
	}

	// Add flags for the deploy command
	// deployCmd.Flags().StringVarP(&brokerType, "broker", "b", "redis", "Specify the broker type (default: redis)")

	// Add deploy command to root
	rootCmd.AddCommand(deployCmd)

	// Execute Cobra root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func deploy() error {
	fmt.Println("deploying...")

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %v", err)
	}
	defer dockerClient.Close()

	fmt.Println("Creating Docker network...")
	err = createNetwork(dockerClient, dockerNetworkName)
	if err != nil {
		return err
	}

	fmt.Println("Starting Redis container...")
	err = startRedis(dockerClient)
	if err != nil {
		return err
	}

	return nil
}

func createNetwork(dockerClien *client.Client, networkName string) error {
	_, err := dockerClien.NetworkInspect(context.Background(), networkName, types.NetworkInspectOptions{})
	if err == nil {
		fmt.Println("Network already exists")
		return nil
	}

	_, err = dockerClien.NetworkCreate(context.Background(), networkName, types.NetworkCreate{})
	if err != nil {
		return fmt.Errorf("error creating network: %v", err)
	}

	fmt.Println("Network created successfully")

	return nil
}

func startRedis(dockerClient *client.Client) error {
	_, err := dockerClient.ContainerInspect(context.Background(), redisContainerName)
	if err == nil {
		fmt.Println("Redis container is already running")
		return nil
	}

	resp, err := dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "redis:latest",
		},
		nil,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				dockerNetworkName: {},
			},
		},
		nil,
		redisContainerName,
	)

	if err != nil {
		return fmt.Errorf("error creating Redis container: %v", err)
	}

	if err := dockerClient.ContainerStart(
		context.Background(),
		resp.ID,
	); err != nil {
		return fmt.Errorf("error starting Redis container: %v", err)
	}

	fmt.Println("Redis container started successfully")

	return nil
}
