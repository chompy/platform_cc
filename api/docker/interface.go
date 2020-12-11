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

package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

// Client defines methods used to interact with Docker.
type Client interface {
	CreateNetwork() error
	DeleteNetwork() error
	GetNetworkHostIP() (string, error)
	CreateNFSVolume(pid string, name string, containerType ObjectContainerType) error
	GetProjectVolumes(pid string) (volume.VolumesListOKBody, error)
	GetAllVolumes() (volume.VolumesListOKBody, error)
	DeleteProjectVolumes(pid string) error
	DeleteAllVolumes() error
	StartContainer(c ContainerConfig) error
	GetProjectContainers(pid string) ([]types.Container, error)
	GetAllContainers() ([]types.Container, error)
	DeleteProjectContainers(pid string) error
	DeleteAllContainers() error
	RunContainerCommand(id string, user string, cmd []string, out io.Writer) error
	UploadDataToContainer(id string, path string, r io.Reader) error
	GetContainerIP(id string) (string, error)
	PullImage(c ContainerConfig) error
	PullImages(containerConfigs []ContainerConfig) error
	ShellContainer(id string, user string, command []string) error
}
