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

import "fmt"

// Config contains configuration for a Docker container.
type Config struct {
	ProjectID    string
	Slot         int  // set the slot to use for mount volumes
	NoCommit     bool // when true don't use committed app image
	ObjectType   ObjectContainerType
	ObjectName   string
	Command      []string
	Image        string
	Volumes      map[string]string
	Binds        map[string]string
	Env          map[string]string
	Ports        []string
	WorkingDir   string
	EnableOSXNFS bool
}

// GetContainerName return the name of the Docker container.
func (d Config) GetContainerName() string {
	return containerName(d.ProjectID, d.ObjectType, d.ObjectName)
}

// GetHumanName returns human readable container name.
func (d Config) GetHumanName() string {
	if d.ObjectType == ObjectContainerRouter {
		return "router"
	}
	return fmt.Sprintf(
		"%s/%s",
		d.ObjectType.TypeName(),
		d.ObjectName,
	)
}

// GetNetworkName returns the name of the Docker network.
func (d Config) GetNetworkName() string {
	return fmt.Sprintf(containerNetworkNameFormat, d.ProjectID)
}

// GetCommand returns the container command.
func (d Config) GetCommand() []string {
	if len(d.Command) > 0 {
		return d.Command
	}
	return []string{"tail", "-f", "/dev/null"}
}

// GetEnv converts environment vars to format needed to start docker container.
func (d Config) GetEnv() []string {
	out := make([]string, 0)
	for k, v := range d.Env {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}
	return out
}
