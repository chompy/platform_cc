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
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// containerVolumeNameFormat is the format for mount volume names.
const containerVolumeNameFormat = dockerNamingPrefix + "v-%s"

// CreateNFSVolume creates a NFS Docker volume.
func (d *Client) CreateNFSVolume(pid string, name string, containerType ObjectContainerType) error {
	pathString := fmt.Sprintf(":/System/Volumes/Data/%s", GetVolumeName(pid, name, containerType))
	_, err := d.cli.VolumeCreate(
		context.Background(),
		volume.VolumesCreateBody{
			Name:   GetVolumeName(pid, name, containerType),
			Driver: "local",
			DriverOpts: map[string]string{
				"type":   "nfs",
				"o":      "addr=host.docker.internal,rw,nolock,hard,nointr,nfsvers=3",
				"device": pathString,
			},
		},
	)
	return tracerr.Wrap(err)
}

// GetProjectVolumes gets a list of all volumes for given project.
func (d *Client) GetProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(dockerNamingPrefix+"*", pid))
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// GetAllVolumes gets a list of all volumes used by PCC.
func (d *Client) GetAllVolumes() (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// DeleteProjectVolumes deletes all volumes for given project.
func (d *Client) DeleteProjectVolumes(pid string) error {
	volList, err := d.GetProjectVolumes(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteVolumes(volList)
}

// DeleteAllVolumes deletes all volumes related to PCC.
func (d *Client) DeleteAllVolumes() error {
	volList, err := d.GetAllVolumes()
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteVolumes(volList)
}

func (d *Client) deleteVolumes(volList volume.VolumesListOKBody) error {
	// prepare progress output
	msgs := make([]string, len(volList.Volumes))
	for i, vol := range volList.Volumes {
		msgs[i] = vol.Name
	}
	prog := output.Progress(msgs)
	// delete volumes
	var wg sync.WaitGroup
	for i, vol := range volList.Volumes {
		wg.Add(1)
		go func(volName string, i int) {
			defer wg.Done()
			if err := d.cli.VolumeRemove(
				context.Background(),
				volName,
				true,
			); err != nil {
				prog(i, output.ProgressMessageError)
				output.Warn(err.Error())
				return
			}
			prog(i, output.ProgressMessageDone)
		}(vol.Name, i)
	}
	wg.Wait()
	return nil
}

// GetVolumeName generates a volume name for given project id and container name.
func GetVolumeName(pid string, name string, containerType ObjectContainerType) string {
	return fmt.Sprintf(dockerNamingPrefix+"%s-%s", pid, string(containerType), name)
}
