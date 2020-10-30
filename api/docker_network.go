package api

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CreateProjectNetwork - create docker network for given project
func (d *dockerClient) CreateProjectNetwork(pid string) error {
	log.Printf("Create Docker network for project '%s.'", pid)
	c := dockerContainerConfig{projectID: pid}
	if _, err := d.cli.NetworkCreate(
		context.Background(),
		c.GetNetworkName(),
		types.NetworkCreate{
			CheckDuplicate: true,
		},
	); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			return err
		}
	}
	return nil
}

// DeleteProjectNetwork - delete docker network for given project
func (d *dockerClient) DeleteProjectNetwork(pid string) error {
	log.Printf("Delete Docker network for project '%s.'", pid)
	c := dockerContainerConfig{projectID: pid}
	err := d.cli.NetworkRemove(
		context.Background(),
		c.GetNetworkName(),
	)
	if !client.IsErrNetworkNotFound(err) {
		return err
	}
	return nil
}
