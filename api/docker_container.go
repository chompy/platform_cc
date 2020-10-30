package api

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

// platformShDockerImagePrefix - prefix for docker images
const platformShDockerImagePrefix = "docker.registry.platform.sh/"

// StartContainer - create and start docker container
func (d *dockerClient) StartContainer(c dockerContainerConfig) error {
	log.Printf("Start Docker container for %s '%s'", c.objectType.TypeName(), c.objectName)
	// get mounts
	mounts := make([]mount.Mount, 0)
	for k, v := range c.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: k,
			Target: v,
		})
	}
	for k, v := range c.Binds {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: k,
			Target: v,
		})
	}
	// create app container
	resp, err := d.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        c.Image,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          c.GetCommand(),
		},
		&container.HostConfig{
			AutoRemove: true,
			Mounts:     mounts,
		},
		&network.NetworkingConfig{},
		c.GetContainerName(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			return nil
		}
		return err
	}
	for _, w := range resp.Warnings {
		log.Printf("WARN: %s", w)
	}
	// attach container to project network
	err = d.cli.NetworkConnect(
		context.Background(),
		c.GetNetworkName(),
		resp.ID,
		nil,
	)
	if err != nil {
		return err
	}
	// start container
	err = d.cli.ContainerStart(
		context.Background(),
		resp.ID,
		types.ContainerStartOptions{},
	)
	return err
}

// GetProjectContainers - get list of active containers for given project
func (d *dockerClient) GetProjectContainers(pid string) ([]types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(dockerNamingPrefix+"*", pid))
	return d.cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: filterArgs,
		},
	)
}

// DeleteProjectContainers - delete all containers related to a project
func (d *dockerClient) DeleteProjectContainers(pid string) error {
	timeout := 30 * time.Second
	container, err := d.GetProjectContainers(pid)
	if err != nil {
		return err
	}
	for _, c := range container {
		log.Printf("Delete Docker container '%s.'", c.Names[0])
		if err := d.cli.ContainerStop(
			context.Background(),
			c.ID,
			&timeout,
		); err != nil {
			return err
		}
	}
	return nil
}
