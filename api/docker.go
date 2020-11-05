package api

import (
	"github.com/docker/docker/client"
)

// dockerClient - docker client for platform_cc
type dockerClient struct {
	cli *client.Client
}

// newDockerClient - create docker client
func newDockerClient() (dockerClient, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return dockerClient{}, err
	}
	return dockerClient{
		cli: cli,
	}, nil
}
