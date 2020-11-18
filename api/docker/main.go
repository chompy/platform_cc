package docker

import (
	"github.com/docker/docker/client"
	"github.com/ztrue/tracerr"
)

// Client is the Docker client for PCC.
type Client struct {
	cli *client.Client
}

// NewClient creates a docker client.
func NewClient() (Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return Client{}, tracerr.Wrap(err)
	}
	return Client{
		cli: cli,
	}, nil
}
