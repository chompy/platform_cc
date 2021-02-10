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

// ContainerStatus defines the status of a container.
type ContainerStatus struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Container Container `json:"-"`
	Running   bool      `json:"running"`
	IPAddress string    `json:"ip_address"`
	Slot      int       `json:"slot"`
}

// Status returns list of container status
func (p *Project) Status() []ContainerStatus {
	out := make([]ContainerStatus, 0)
	for _, app := range p.Apps {
		c := p.NewContainer(app)
		status, _ := c.containerHandler.ContainerStatus(c.Config.GetContainerName())
		out = append(out, ContainerStatus{
			Name:      app.Name,
			Type:      app.Type,
			Container: c,
			Running:   status.Running,
			Slot:      status.Slot,
			IPAddress: status.IPAddress,
		})
	}
	for _, service := range p.Services {
		c := p.NewContainer(service)
		status, _ := c.containerHandler.ContainerStatus(c.Config.GetContainerName())
		out = append(out, ContainerStatus{
			Name:      service.Name,
			Type:      service.Type,
			Container: c,
			Running:   status.Running,
			Slot:      status.Slot,
			IPAddress: status.IPAddress,
		})
	}
	return out
}
