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
