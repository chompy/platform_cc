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

package project

import (
	"runtime"
	"strings"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/docker"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

const containerDataDirectory = "/mnt/data"

func isMacOS() bool {
	return runtime.GOOS == "darwin"
}

func prepareNfsVolume(p *Project, d interface{}) error {
	name := ""
	containerType := docker.ObjectContainerApp
	switch d.(type) {
	case def.App:
		{
			name = d.(def.App).Name
			break
		}
	case def.AppWorker:
		{
			name = d.(def.AppWorker).Name
			containerType = docker.ObjectContainerWorker
			break
		}
	}
	if !isMacOS() {
		return nil
	}
	if err := p.docker.CreateNFSVolume(p.ID, name+"-nfs", containerType); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			return tracerr.Wrap(err)
		}
	}
	return nil
}
