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
	"gitlab.com/contextualcode/platform_cc/api/container"
)

// Status returns list of container status
func (p *Project) Status() []container.Status {
	out := make([]container.Status, 0)
	for _, app := range p.Apps {
		c := p.NewContainer(app)
		status, _ := c.containerHandler.ContainerStatus(c.Config.GetContainerName())
		if status.Name == "" {
			status.Name = c.Config.ObjectName
			status.ObjectType = c.Config.ObjectType
			status.Type = app.Type
		}
		out = append(out, status)
		for _, worker := range app.Workers {
			wc := p.NewContainer(worker)
			status, _ := c.containerHandler.ContainerStatus(wc.Config.GetContainerName())
			if status.Name == "" {
				status.Name = wc.Config.ObjectName
				status.ObjectType = wc.Config.ObjectType
				status.Type = worker.Type
			}
			out = append(out, status)
		}
	}
	for _, service := range p.Services {
		c := p.NewContainer(service)
		status, _ := c.containerHandler.ContainerStatus(c.Config.GetContainerName())
		if status.Name == "" {
			status.Name = c.Config.ObjectName
			status.ObjectType = c.Config.ObjectType
			status.Type = service.Type
		}
		out = append(out, status)
	}
	return out
}
