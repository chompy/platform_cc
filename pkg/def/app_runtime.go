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

// AppRuntime defines runtime configuration.
type AppRuntime struct {
	RequestTerminateTimeout int                    `yaml:"request_terminate_timeout" json:"-"`
	Extensions              []*AppRuntimeExtension `yaml:"extensions" json:"extensions"`
	DisabledExtensions      []string               `yaml:"disabled_extensions" json:"disabled_extensions,omitempty"`
	Xdebug                  AppRuntimeXdebug       `yaml:"xdebug" json:"xdebug"`
}

// SetDefaults sets the default values.
func (d *AppRuntime) SetDefaults() {
	if d.RequestTerminateTimeout <= 0 {
		d.RequestTerminateTimeout = 300
	}
	if d.Extensions == nil || len(d.Extensions) == 0 {
		d.Extensions = make([]*AppRuntimeExtension, 0)
	}
	for i := range d.Extensions {
		d.Extensions[i].SetDefaults()
	}
}

// Validate checks for errors.
func (d AppRuntime) Validate(root *App) []error {
	o := make([]error, 0)
	for i := range d.Extensions {
		if e := d.Extensions[i].Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
