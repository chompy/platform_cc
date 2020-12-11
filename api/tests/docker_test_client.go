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

package tests

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"gitlab.com/contextualcode/platform_cc/api/docker"
)

const dockerTestNetworkIP = "192.1.2.3"
const dockerTestContainerIPD = "192.1.2.4"

// MockDockerClient is a test Docker client that makes no real connects.
type MockDockerClient struct {
	hasNetwork    *bool
	volumeList    []string
	containerList []string
}

// CreateNetwork simulates the creation of a Docker network.
func (d MockDockerClient) CreateNetwork() error {
	*d.hasNetwork = true
	return nil
}

// DeleteNetwork simulates the deletion of a Docker network.
func (d MockDockerClient) DeleteNetwork() error {
	*d.hasNetwork = false
	return nil
}

// GetNetworkHostIP returns a fake IP if network we created, otherwise error.
func (d MockDockerClient) GetNetworkHostIP() (string, error) {
	if *d.hasNetwork {
		return dockerTestNetworkIP, nil
	}
	return "", fmt.Errorf("network not started")
}

// CreateNFSVolume simulates the creation of a Docker NFS volume.
func (d MockDockerClient) CreateNFSVolume(pid string, name string, containerType docker.ObjectContainerType) error {
	return nil
}

// GetProjectVolumes returns a list of volumes for given project.
func (d MockDockerClient) GetProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	return volume.VolumesListOKBody{}, nil
}

// GetAllVolumes returns list of all volumes.
func (d MockDockerClient) GetAllVolumes() (volume.VolumesListOKBody, error) {
	return volume.VolumesListOKBody{}, nil
}

// DeleteProjectVolumes simulates deletion of all project volumes.
func (d MockDockerClient) DeleteProjectVolumes(pid string) error {
	return nil
}

// DeleteAllVolumes simulates deletion of all Docker volumes.
func (d MockDockerClient) DeleteAllVolumes() error {
	return nil
}

// StartContainer simulates starting a Docker container.
func (d MockDockerClient) StartContainer(c docker.ContainerConfig) error {
	for _, name := range d.containerList {
		if name == c.GetContainerName() {
			return nil
		}
	}
	d.containerList = append(d.containerList, c.GetContainerName())
	return nil
}

// GetProjectContainers returns list of project Docker containers.
func (d MockDockerClient) GetProjectContainers(pid string) ([]types.Container, error) {
	return []types.Container{}, nil
}

// GetAllContainers returns list of all Docker containers.
func (d MockDockerClient) GetAllContainers() ([]types.Container, error) {
	return nil, nil
}

// DeleteProjectContainers simulates deletion of project Docker containers.
func (d MockDockerClient) DeleteProjectContainers(pid string) error {
	return nil
}

// DeleteAllContainers simulates deletion of all Docker containers.
func (d MockDockerClient) DeleteAllContainers() error {
	return nil
}

// RunContainerCommand simulates running a Docker container command.
func (d MockDockerClient) RunContainerCommand(id string, user string, cmd []string, out io.Writer) error {
	return nil
}

// UploadDataToContainer simulates uploading data to a Docker container.
func (d MockDockerClient) UploadDataToContainer(id string, path string, r io.Reader) error {
	return nil
}

// GetContainerIP returns a fake container IP if fake container was started.
func (d MockDockerClient) GetContainerIP(id string) (string, error) {
	return dockerTestNetworkIP, nil
}

// PullImage simulates pulling a Docker image.
func (d MockDockerClient) PullImage(c docker.ContainerConfig) error {
	return nil
}

// PullImages simulates pulling multiple Docker images.
func (d MockDockerClient) PullImages(containerConfigs []docker.ContainerConfig) error {
	return nil
}

// ShellContainer simulates shelling in to a Docker container.
func (d MockDockerClient) ShellContainer(id string, user string, command []string) error {
	return nil
}
