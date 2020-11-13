package api

import (
	"github.com/docker/docker/client"
)

// DockerClient - docker client for platform_cc
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient - create docker client
func NewDockerClient() (DockerClient, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return DockerClient{}, err
	}
	return DockerClient{
		cli: cli,
	}, nil
}
