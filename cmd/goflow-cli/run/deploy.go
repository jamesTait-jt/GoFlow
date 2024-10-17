package run

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/jamesTait-jt/goflow/cmd/goflow-cli/docker"
)

var dockerNetworkName = "goflow-network"

var redisContainerName = "redis-server"

var redisImage = "redis:server"
var pluginBuilderImage = "plugin-builder"
var workerpoolImage = "workerpool"

func Deploy() error {
	dockerClient, err := docker.New()
	if err != nil {
		return fmt.Errorf("error creating Docker client: %v", err)
	}
	defer dockerClient.Close()

	fmt.Println("Creating Docker network...")
	err = dockerClient.CreateNetwork(dockerNetworkName)
	if err != nil {
		return err
	}

	fmt.Println("Starting Redis container...")
	err = startRedis(dockerClient)
	if err != nil {
		return err
	}

	fmt.Println("Compiling plugins...")
	err = compilePlugins(dockerClient)
	if err != nil {
		return err
	}

	fmt.Println("Starting WorkerPool container...")
	err = startWorkerPool(dockerClient)
	if err != nil {
		return err
	}

	fmt.Println("Deployment successful!")
	return nil
}

func startRedis(dockerClient *docker.Docker) error {
	exists, running, containerID, err := dockerClient.ContainerInfo(redisContainerName)
	if err != nil {
		return err
	}

	if running {
		fmt.Println("Redis container already started")

		return nil
	}

	if !exists {
		err = dockerClient.PullImage(redisContainerName)
		if err != nil {
			return fmt.Errorf("failed to pull redis image: %v", err)
		}

		containerID, err = dockerClient.CreateContainer(
			&container.Config{
				Image: "redis:latest",
			},
			&container.HostConfig{
				PortBindings: nat.PortMap{
					"6379/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0", // Listen on all network interfaces
							HostPort: "6379",    // Expose on this port on the host
						},
					},
				},
			},
			dockerNetworkName,
			redisContainerName,
		)

		if err != nil {
			return fmt.Errorf("error creating Redis container: %v", err)
		}
	}

	if err = dockerClient.StartContainer(containerID); err != nil {
		return fmt.Errorf("error starting redis container: %v", err)
	}

	fmt.Println("Redis container started successfully")

	return nil
}

func compilePlugins(dockerClient *docker.Docker) error {
	// TODO: Make this better
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	handlersPath := fmt.Sprintf("%s/handlers", cwd)

	containerID, err := dockerClient.CreateContainer(
		&container.Config{
			Image: pluginBuilderImage,
			Cmd:   []string{"handlers"},
		},
		&container.HostConfig{
			Binds:      []string{fmt.Sprintf("%s:/app/handlers", handlersPath)},
			AutoRemove: true,
		},
		"",
		"",
	)

	if err != nil {
		return fmt.Errorf("failed to create plugin-builder container: %v", err)
	}

	if err = dockerClient.StartContainer(containerID); err != nil {
		return fmt.Errorf("failed to start plugin-builder container: %v", err)
	}

	dockerClient.WaitForContainerToFinish(containerID)

	containerPassed, err := dockerClient.ContainerPassed(containerID)
	if err != nil {
		return fmt.Errorf("failed to check plugin-builder exit status: %v", err)
	}

	if !containerPassed {
		return errors.New("plugin-builder container failed")
	}

	fmt.Println("plugins compiled!")

	return nil
}

func startWorkerPool(dockerClient *docker.Docker) error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	handlersPath := fmt.Sprintf("%s/handlers", cwd)

	hostConfig := &container.HostConfig{
		Binds:      []string{fmt.Sprintf("%s:/app/handlers", handlersPath)},
		AutoRemove: true,
	}

	containerID, err := dockerClient.CreateContainer(
		&container.Config{
			Image: workerpoolImage,
			Cmd: []string{
				"--broker-type", "redis",
				"--broker-addr", fmt.Sprintf("%s:6379", redisContainerName),
				"--handlers-path", "/app/handlers/compiled",
			},
		},
		hostConfig,
		dockerNetworkName,
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to create workerpool container: %v", err)
	}

	if err := dockerClient.StartContainer(containerID); err != nil {
		return fmt.Errorf("error starting Redis container: %v", err)
	}

	fmt.Println("WorkerPool container started successfully")

	return nil
}
