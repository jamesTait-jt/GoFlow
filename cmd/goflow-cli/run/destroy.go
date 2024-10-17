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

	for _, containerID := range []string{config.RedisContainerName, config.WorkerpoolContainerName} {
		fmt.Printf("Destroying container '%s'\n", containerID)

		if err = dockerClient.DestroyContainer(containerID); err != nil {
			return fmt.Errorf("failed to destroy container '%s': %v", containerID, err)
		}
	}

	fmt.Println("Destroying Docker network...")

	if err = dockerClient.DestroyNetwork(config.DockerNetworkID); err != nil {
		return fmt.Errorf("failed to destroy network '%s': %v", config.DockerNetworkID, err)
	}

	fmt.Println("Done!")

	return nil
}
