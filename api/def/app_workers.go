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

package def

// AppWorker defines a worker.
type AppWorker struct {
	Size          string                            `yaml:"size"`
	Disk          int                               `yaml:"disk"`
	Mounts        map[string]*AppMount              `yaml:"mounts"`
	Relationships map[string]string                 `yaml:"relationships"`
	Variables     map[string]map[string]interface{} `yaml:"variables"`
}

// SetDefaults sets the default values.
func (d *AppWorker) SetDefaults() {
	for k := range d.Mounts {
		d.Mounts[k].SetDefaults()
	}
	if d.Size == "" {
		d.Size = "S"
	}
	if d.Disk < 256 {
		d.Disk = 256
	}
}

// Validate checks for errors.
func (d AppWorker) Validate() []error {
	// TODO
	return []error{}
}
