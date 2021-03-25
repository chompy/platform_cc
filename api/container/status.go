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

// Status defines container status.
type Status struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	ObjectType   ObjectContainerType `json:"object_type"`
	Image        string              `json:"image"`
	Type         string              `json:"type"`
	ProjectID    string              `json:"project_id"`
	Committed    bool                `json:"committed"`
	Running      bool                `json:"running"`
	State        string              `json:"state"`
	IPAddress    string              `json:"ip_address"`
	Slot         int                 `json:"slot"`
	HasContainer bool                `json:"has_container"`
}
