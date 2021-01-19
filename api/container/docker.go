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
	"github.com/docker/docker/client"
	"github.com/ztrue/tracerr"
)

// Docker defines the Docker container handler.
type Docker struct {
	client *client.Client
}

// NewDocker creates a new Docker container handler.
func NewDocker() (Docker, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return Docker{}, tracerr.Wrap(err)
	}
	return Docker{
		client: dockerClient,
	}, nil
}
