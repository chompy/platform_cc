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

import "github.com/ztrue/tracerr"

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
	return nil
}
