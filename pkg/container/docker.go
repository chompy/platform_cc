/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package container

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// Docker defines the Docker container handler.
type Docker struct {
	client *client.Client
}

// NewDocker creates a new Docker container handler.
func NewDocker() (Docker, error) {
	dockerOpts, err := getDockerOpts()
	if err != nil {
		return Docker{}, errors.WithStack(err)
	}
	dockerClient, err := client.NewClientWithOpts(dockerOpts...)
	if err != nil {
		return Docker{}, errors.WithStack(convertDockerError(err))
	}
	return Docker{
		client: dockerClient,
	}, nil
}

func getDockerOpts() ([]client.Opt, error) {
	timeout := client.WithTimeout(time.Second * 900) // 15 minutes
	// standard docker environment
	host := os.Getenv("DOCKER_HOST")
	if !strings.HasPrefix(host, "ssh://") {
		return []client.Opt{client.FromEnv, timeout}, nil
	}
	// use ssh
	// https://stackoverflow.com/a/57933792
	helper, err := connhelper.GetConnectionHelper(host)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}
	var clientOpts []client.Opt
	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
		timeout,
	)
	version := os.Getenv("DOCKER_API_VERSION")
	if version != "" {
		clientOpts = append(clientOpts, client.WithVersion(version))
	} else {
		clientOpts = append(clientOpts, client.WithAPIVersionNegotiation())
	}
	return clientOpts, nil
}

// convertDockerError converts internal docker error to platform.cc error.
func convertDockerError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "No such container") {
		// container not found
		return ErrContainerNotFound
	} else if strings.Contains(err.Error(), "No such image") || strings.Contains(err.Error(), "manifest unknown") {
		// imagen not found
		return ErrImageNotFound
	}
	// nothing found, return original error
	return err
}
