package docker

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ztrue/tracerr"
)

// globalNetworkName is the name of the global network.
const globalNetworkName = "pcc"

// CreateNetwork creates a global network for use with all PCC containers.
func (d *Client) CreateNetwork() error {
	log.Println("Create network.")
	if _, err := d.cli.NetworkCreate(
		context.Background(),
		globalNetworkName,
		types.NetworkCreate{
			CheckDuplicate: true,
		},
	); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			return tracerr.Wrap(err)
		}
	}
	return nil
}

// DeleteNetwork deletes the global network.
func (d *Client) DeleteNetwork() error {
	log.Println("Delete network.")
	err := d.cli.NetworkRemove(
		context.Background(),
		globalNetworkName,
	)
	if !client.IsErrNetworkNotFound(err) {
		return tracerr.Wrap(err)
	}
	return nil
}

// GetNetworkHostIP gets the host IP address for the given project's network.
func (d *Client) GetNetworkHostIP() (string, error) {
	net, err := d.cli.NetworkInspect(
		context.Background(),
		globalNetworkName,
	)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return net.IPAM.Config[0].Gateway, nil
}
