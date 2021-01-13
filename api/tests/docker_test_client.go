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
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"gitlab.com/contextualcode/platform_cc/api/docker"
)

const dockerTestNetworkIP = "192.1.2.3"
const dockerTestContainerIPD = "192.1.2.4"

// MockDockerClient is a test Docker client that makes no real connects.
type MockDockerClient struct {
	hasNetwork    *bool
	volumeList    *[]string
	containerList *[]string
}

// NewMockDockerClient creates a new MockDockerClient.
func NewMockDockerClient() MockDockerClient {
	hasNetwork := false
	volumeList := make([]string, 0)
	containerList := make([]string, 0)
	return MockDockerClient{
		hasNetwork:    &hasNetwork,
		volumeList:    &volumeList,
		containerList: &containerList,
	}
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
	if d.hasNetwork != nil && *d.hasNetwork {
		return dockerTestNetworkIP, nil
	}
	return "", fmt.Errorf("network not started")
}

// CreateNFSVolume simulates the creation of a Docker NFS volume.
func (d MockDockerClient) CreateNFSVolume(pid string, name string, path string, containerType docker.ObjectContainerType) error {
	return nil
}

// GetProjectVolumes returns a list of volumes for given project.
func (d MockDockerClient) GetProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	out := make([]*types.Volume, 0)
	for _, name := range *d.volumeList {
		if strings.Contains(name, pid) {
			out = append(out, &types.Volume{
				Name: name,
			})
		}
	}
	return volume.VolumesListOKBody{
		Volumes: out,
	}, nil
}

// GetAllVolumes returns list of all volumes.
func (d MockDockerClient) GetAllVolumes() (volume.VolumesListOKBody, error) {
	out := make([]*types.Volume, 0)
	for _, name := range *d.volumeList {
		out = append(out, &types.Volume{
			Name: name,
		})
	}
	return volume.VolumesListOKBody{
		Volumes: out,
	}, nil
}

// DeleteProjectVolumes simulates deletion of all project volumes.
func (d MockDockerClient) DeleteProjectVolumes(pid string) error {
	out := make([]string, 0)
	for _, name := range *d.volumeList {
		if !strings.Contains(name, pid) {
			out = append(out, name)
		}
	}
	*d.volumeList = out
	return nil
}

// DeleteAllVolumes simulates deletion of all Docker volumes.
func (d MockDockerClient) DeleteAllVolumes() error {
	*d.volumeList = make([]string, 0)
	return nil
}

// StartContainer simulates starting a Docker container.
func (d MockDockerClient) StartContainer(c docker.ContainerConfig) error {
	for _, name := range *d.containerList {
		if name == c.GetContainerName() {
			return nil
		}
	}
	*d.containerList = append(*d.containerList, c.GetContainerName())
	for name := range c.Volumes {
		*d.volumeList = append(*d.volumeList, name)
	}
	return nil
}

// GetProjectContainers returns list of project Docker containers.
func (d MockDockerClient) GetProjectContainers(pid string) ([]types.Container, error) {
	out := make([]types.Container, 0)
	for _, name := range *d.containerList {
		if strings.Contains(name, pid) {
			out = append(out, types.Container{
				ID:    name,
				Names: []string{name},
			})
		}
	}
	return out, nil
}

// GetAllContainers returns list of all Docker containers.
func (d MockDockerClient) GetAllContainers() ([]types.Container, error) {
	out := make([]types.Container, 0)
	for _, name := range *d.containerList {
		out = append(out, types.Container{
			ID:    name,
			Names: []string{name},
		})
	}
	return out, nil
}

// DeleteProjectContainers simulates deletion of project Docker containers.
func (d MockDockerClient) DeleteProjectContainers(pid string) error {
	out := make([]string, 0)
	for _, name := range *d.containerList {
		if !strings.Contains(name, pid) {
			out = append(out, name)
		}
	}
	*d.containerList = out
	return nil
}

// DeleteAllContainers simulates deletion of all Docker containers.
func (d MockDockerClient) DeleteAllContainers() error {
	*d.containerList = make([]string, 0)
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
func (d MockDockerClient) ShellContainer(id string, user string, command []string, stdin io.Reader) error {
	return nil
}

// ContainerLog simulates container logging.
func (d MockDockerClient) ContainerLog(id string) {
}

// CommitContainer simulates commiting a container.
func (d MockDockerClient) CommitContainer(id string) error {
	return nil
}

// DeleteProjectImages simulates deleting project images.
func (d MockDockerClient) DeleteProjectImages(pid string) error {
	return nil
}

// DeleteAllImages simulates deleting all PCC images.
func (d MockDockerClient) DeleteAllImages() error {
	return nil
}

// GetCommittedImage simulates retrivial of committed container image.
func (d MockDockerClient) GetCommittedImage(c docker.ContainerConfig) (string, error) {
	return "", nil
}
