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
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// createNFSVolume creates a NFS Docker volume.
func (d Docker) createNFSVolume(pid string, name string, path string, containerType ObjectContainerType) error {
	path = strings.TrimLeft(path, "/")
	pathString := fmt.Sprintf(":/System/Volumes/Data/%s", path)
	output.LogDebug("NFS mount path.", pathString)
	volCreateBody := volume.VolumesCreateBody{
		Name:   getMountName(pid, name, containerType),
		Driver: "local",
		DriverOpts: map[string]string{
			"type":   "nfs",
			"o":      "addr=host.docker.internal,rw,nolock,hard,nointr,nfsvers=3",
			"device": pathString,
		},
	}
	output.LogDebug("Create Docker NFS volume.", volCreateBody)
	v, err := d.client.VolumeCreate(
		context.Background(),
		volCreateBody,
	)
	output.LogDebug("Created Docker NFS volume.", v)
	return tracerr.Wrap(err)
}

// listProjectVolumes gets a list of all volumes for given project.
func (d Docker) listProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	output.LogDebug(fmt.Sprintf("List volumes for project %s (all slots).", pid), nil)
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(containerNamingPrefix+"*", pid))
	return d.client.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// listProjectSlotVolumes gets a list of all volumes for given project slot.
func (d Docker) listProjectSlotVolumes(pid string, slot int) (volume.VolumesListOKBody, error) {
	output.LogDebug(fmt.Sprintf("List volumes for project %s in slot %d.", pid, slot), nil)
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(containerNamingPrefix+"*", pid))
	list, err := d.client.VolumeList(
		context.Background(),
		filterArgs,
	)
	if err != nil {
		return volume.VolumesListOKBody{}, tracerr.Wrap(err)
	}
	out := make([]*types.Volume, 0)
	for _, v := range list.Volumes {
		if volumeBelongsToSlot(v.Name, slot) {
			out = append(out, v)
		}
	}
	list.Volumes = out
	return list, nil
}

// listAllVolumes gets a list of all volumes used by Platform.CC.
func (d Docker) listAllVolumes() (volume.VolumesListOKBody, error) {
	output.LogDebug("List all volumes.", nil)
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	return d.client.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// deleteVolumes deletes given Docker volumes.
func (d Docker) deleteVolumes(volList volume.VolumesListOKBody) error {
	// prepare progress output
	output.LogDebug("Delete Docker volumes.", volList)
	msgs := make([]string, len(volList.Volumes))
	for i, vol := range volList.Volumes {
		msgs[i] = vol.Name
	}
	done := output.Duration("Delete volumes.")
	prog := output.Progress(msgs)
	// delete volumes
	var wg sync.WaitGroup
	for i, vol := range volList.Volumes {
		wg.Add(1)
		go func(volName string, i int) {
			defer wg.Done()
			if err := d.client.VolumeRemove(
				context.Background(),
				volName,
				true,
			); err != nil {
				prog(i, output.ProgressMessageError, nil, nil)
				output.Warn(err.Error())
				return
			}
			prog(i, output.ProgressMessageDone, nil, nil)
		}(vol.Name, i)
	}
	wg.Wait()
	done()
	return nil
}
