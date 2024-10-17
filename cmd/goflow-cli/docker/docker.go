package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Docker struct {
	ctx    context.Context
	client *client.Client
}

func New() (*Docker, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Docker{
		ctx:    context.Background(),
		client: c,
	}, nil
}

func (d *Docker) CreateNetwork(networkID string) error {
	_, err := d.client.NetworkInspect(d.ctx, networkID, network.InspectOptions{})
	if err == nil {
		fmt.Println("Network already exists")
		return nil
	}

	_, err = d.client.NetworkCreate(d.ctx, networkID, network.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating network: %v", err)
	}

	fmt.Println("Network created successfully")

	return nil
}

func (d *Docker) DestroyNetwork(networkID string) error {
	return d.client.NetworkRemove(d.ctx, networkID)
}

type ContainerInfo struct {
	Exists  bool
	Running bool
	ID      string
}

func (d *Docker) ContainerInfo(containerName string) (ContainerInfo, error) {
	containerJSON, err := d.client.ContainerInspect(d.ctx, containerName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return ContainerInfo{false, false, ""}, nil
		}

		return ContainerInfo{}, fmt.Errorf("failed to inspect container '%s': %v", containerName, err)
	}

	if containerJSON.State.Running {
		return ContainerInfo{true, true, containerJSON.ID}, nil
	}

	return ContainerInfo{true, false, containerJSON.ID}, nil
}

func (d *Docker) CreateContainer(
	config *container.Config,
	hostConfig *container.HostConfig,
	networkID string,
	containerName string,
) (string, error) {
	var networkConfig *network.NetworkingConfig
	if networkID != "" {
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkID: {},
			},
		}
	}

	resp, err := d.client.ContainerCreate(
		d.ctx,
		config,
		hostConfig,
		networkConfig,
		nil,
		containerName,
	)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (d *Docker) StartContainer(containerID string) error {
	if err := d.client.ContainerStart(
		d.ctx, containerID, container.StartOptions{},
	); err != nil {
		return err
	}

	return nil
}

func (d *Docker) WaitForContainerToFinish(containerID string) error {
	statusCh, errCh := d.client.ContainerWait(d.ctx, containerID, container.WaitConditionNotRunning)

	select {
	case <-statusCh:
		return nil

	case err := <-errCh:
		return err
	}
}

func (d *Docker) ContainerPassed(containerID string) (bool, error) {
	containerInspect, err := d.client.ContainerInspect(d.ctx, containerID)
	if err != nil {
		return false, err
	}

	return containerInspect.State.ExitCode == 0, nil
}

func (d *Docker) DestroyContainer(containerID string) error {
	if err := d.client.ContainerStop(d.ctx, containerID, container.StopOptions{}); err != nil {
		return err
	}

	return d.client.ContainerRemove(
		d.ctx,
		containerID,
		container.RemoveOptions{RemoveVolumes: true, Force: true},
	)
}

func (d *Docker) ImageExistsLocally(imageTag string) (bool, error) {
	images, err := d.client.ImageList(d.ctx, image.ListOptions{})
	if err != nil {
		return false, err
	}

	for i := 0; i < len(images); i++ {
		for _, tag := range images[i].RepoTags {
			if tag == imageTag {
				return true, nil
			}
		}
	}

	return false, nil
}

func (d *Docker) PullImage(imageTag string) error {
	imageExists, err := d.ImageExistsLocally(imageTag)
	if err != nil {
		return err
	}

	if imageExists {
		return nil
	}

	fmt.Printf("Pulling %s...\n", imageTag)

	resp, err := d.client.ImagePull(d.ctx, imageTag, image.PullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	_, err = io.ReadAll(resp)
	if err != nil {
		return err
	}

	return err
}

func (d *Docker) Close() {
	d.client.Close()
}
