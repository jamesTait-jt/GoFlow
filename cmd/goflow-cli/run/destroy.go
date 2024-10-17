package run

import (
	"fmt"

	"github.com/jamesTait-jt/goflow/cmd/goflow-cli/config"
	"github.com/jamesTait-jt/goflow/cmd/goflow-cli/docker"
)

func Destroy() error {
	dockerClient, err := docker.New()
	if err != nil {
		return fmt.Errorf("error creating Docker client: %v", err)
	}
	defer dockerClient.Close()

	fmt.Printf("Destroying container '%s'...\n", config.RedisContainerName)

	if err = dockerClient.DestroyContainer(config.RedisContainerName); err != nil {
		return fmt.Errorf("failed to destroy container '%s': %v", config.RedisContainerName, err)
	}

	fmt.Printf("Destroying container '%s'...\n", config.WorkerpoolContainerName)

	if err = dockerClient.DestroyContainer(config.WorkerpoolContainerName); err != nil {
		return fmt.Errorf("failed to destroy container '%s': %v", config.WorkerpoolContainerName, err)
	}

	fmt.Println("Destroying Docker network...")

	if err = dockerClient.DestroyNetwork(config.DockerNetworkID); err != nil {
		return fmt.Errorf("failed to destroy network '%s': %v", config.DockerNetworkID, err)
	}

	return nil
}
