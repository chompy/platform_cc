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
	Name          string                            `json:"-"`
	Path          string                            `json:"-"`
	Type          string                            `json:"-"`
	Runtime       AppRuntime                        `json:"-"`
	Dependencies  AppDependencies                   `json:"-"`
	Size          string                            `yaml:"size" json:"size"`
	Disk          int                               `yaml:"disk" json:"disk"`
	Mounts        map[string]*AppMount              `yaml:"mounts" json:"mounts"`
	Relationships map[string]string                 `yaml:"relationships" json:"relationship"`
	Variables     map[string]map[string]interface{} `yaml:"variables" json:"variables"`
	Commands      AppWorkersCommands                `yaml:"commands" json:"commands"`
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
func (d AppWorker) Validate(root *App) []error {
	o := make([]error, 0)
	if e := d.Commands.Validate(root); len(e) > 0 {
		o = append(o, e...)
	}
	for _, mount := range d.Mounts {
		if e := mount.Validate(root); len(e) > 0 {
			o = append(o, e...)
		}
	}
	return o
}
