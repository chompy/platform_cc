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

// Package docker provides an abstraction to the Docker API.
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
