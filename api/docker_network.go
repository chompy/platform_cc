package api

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// CreateProjectNetwork - create docker network for given project
func (d *DockerClient) CreateProjectNetwork(pid string) error {
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
func (d *DockerClient) DeleteProjectNetwork(pid string) error {
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

// DeleteAllNetworks - delete docker networks related to PCC
func (d *DockerClient) DeleteAllNetworks() error {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	networks, err := d.cli.NetworkList(
		context.Background(),
		types.NetworkListOptions{
			Filters: filterArgs,
		},
	)
	if err != nil {
		return err
	}
	for _, network := range networks {
		log.Printf("Delete Docker network '%s.'", network.Name)
		if err := d.cli.NetworkRemove(
			context.Background(),
			network.ID,
		); err != nil {
			return err
		}
	}
	return nil
}

// GetNetworkHostIP - get host ip address for network
func (d *DockerClient) GetNetworkHostIP(pid string) (string, error) {
	c := dockerContainerConfig{projectID: pid}
	net, err := d.cli.NetworkInspect(
		context.Background(),
		c.GetNetworkName(),
	)
	if err != nil {
		return "", err
	}
	return net.IPAM.Config[0].Gateway, nil
}
