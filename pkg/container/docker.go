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
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// Docker defines the Docker container handler.
type Docker struct {
	client *client.Client
}

// NewDocker creates a new Docker container handler.
func NewDocker() (Docker, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return Docker{}, errors.WithStack(convertDockerError(err))
	}
	return Docker{
		client: dockerClient,
	}, nil
}

// convertDockerError converts internal docker error to platform.cc error.
func convertDockerError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "No such container") {
		// container not found
		return ErrContainerNotFound
	} else if strings.Contains(err.Error(), "No such image") {
		// imagen not found
		return ErrImageNotFound
	}
	// nothing found, return original error
	return err
}
