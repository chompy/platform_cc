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
	"github.com/docker/docker/api/types/mount"
	"github.com/ztrue/tracerr"
)

// AllStop stops all Platform.CC Docker containers.
func (d Docker) AllStop() error {
	containers, err := d.listAllContainers()
	if err != nil {
		return tracerr.Wrap(err)
	}
	return tracerr.Wrap(d.deleteContainers(containers))
}

// AllPurge deletes all Platform.CC Docker resources.
func (d Docker) AllPurge() error {
	// stop
	if err := d.AllStop(); err != nil {
		return tracerr.Wrap(err)
	}
	// delete volumes
	volList, err := d.listAllVolumes()
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := d.deleteVolumes(volList); err != nil {
		return tracerr.Wrap(err)
	}
	// delete images
	imgList, err := d.listAllImages()
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := d.deleteImages(imgList); err != nil {
		return tracerr.Wrap(err)
	}
	// delete network
	if err := d.deleteNetwork(); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// AllStatus returns the status of all Platform.CC Docker containers.
func (d Docker) AllStatus() ([]Status, error) {
	containers, err := d.listAllContainers()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	out := make([]Status, len(containers))
	for i, c := range containers {
		// get ip address
		ipAddress := ""
		for name, network := range c.NetworkSettings.Networks {
			if name == dockerNetworkName {
				ipAddress = network.IPAddress
				break
			}
		}
		// get slot
		slot := 1
		for _, m := range c.Mounts {
			if m.Type == mount.TypeVolume {
				slot = volumeGetSlot(m.Name)
				break
			}
		}
		// get service name
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		config := containerConfigFromName(name)
		out[i] = Status{
			ID:         c.ID,
			Name:       config.ObjectName,
			ObjectType: config.ObjectType,
			Image:      c.Image,
			Type:       d.serviceTypeFromImage(c.Image),
			Committed:  d.hasCommit(c.Names[0]),
			ProjectID:  config.ProjectID,
			Running:    c.State == "running",
			State:      c.State,
			IPAddress:  ipAddress,
			Slot:       slot,
		}
	}
	return out, nil
}
