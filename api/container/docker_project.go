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
	"log"

	"github.com/ztrue/tracerr"
)

// ProjectStop stops all running Docker containers for given project.
func (d Docker) ProjectStop(pid string) error {
	containers, err := d.listProjectContainers(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteContainers(containers)
}

// ProjectPurge deletes all Docker resources for given project.
func (d Docker) ProjectPurge(pid string) error {
	// stop
	if err := d.ProjectStop(pid); err != nil {
		return tracerr.Wrap(err)
	}
	// delete volumes
	vols, err := d.listProjectVolumes(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := d.deleteVolumes(vols); err != nil {
		return tracerr.Wrap(err)
	}
	// delete images
	images, err := d.listProjectImages(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := d.deleteImages(images); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// ProjectPurgeSlot deletes all Docker resources for given project slot.
func (d Docker) ProjectPurgeSlot(pid string, slot int) error {
	// stop
	if err := d.ProjectStop(pid); err != nil {
		return tracerr.Wrap(err)
	}
	// delete volumes
	vols, err := d.listProjectSlotVolumes(pid, slot)
	if err != nil {
		log.Println(err)
		return tracerr.Wrap(err)
	}
	if err := d.deleteVolumes(vols); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// ProjectCopySlot copies volumes in given slot to another slot.
/*func (d Docker) ProjectCopySlot(pid string, sourceSlot int, destSlot int) error {
	// can't have same slots
	if sourceSlot == destSlot {
		return tracerr.New("source and destination slots cannot be the same")
	}
	// purge dest slot in prep for copy
	d.ProjectPurgeSlot(pid, destSlot)

	// get list of source slot volumes
	volList, err := d.listProjectSlotVolumes(pid, sourceSlot)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// create container
	cConfig := &container.Config{
		Image: "busybox",
		Cmd:   []string{""},
	}
	cHostConfig := &container.HostConfig{
		AutoRemove:   true,
		Privileged:   true,
		Mounts:       mounts,
		PortBindings: portBinding,
	}
	output.LogDebug(fmt.Sprintf("Container create. (Name %s)", c.GetContainerName()), []interface{}{cConfig, cHostConfig})
	resp, err := d.client.ContainerCreate(
		context.Background(),
		cConfig,
		cHostConfig,
		&network.NetworkingConfig{},
		c.GetContainerName(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			return nil
		}
		return tracerr.Wrap(err)
	}
	output.LogDebug("Container created.", resp)
	// attach container to project network
	if err := d.client.NetworkConnect(
		context.Background(),
		dockerNetworkName,
		resp.ID,
		nil,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// start container
	if err := d.client.ContainerStart(
		context.Background(),
		resp.ID,
		types.ContainerStartOptions{},
	); err != nil {
		return tracerr.Wrap(err)
	}

}*/
