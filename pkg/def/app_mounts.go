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

import (
	"strings"
)

// AppMount defines persistent mount volumes
type AppMount struct {
	Source     string `yaml:"source" json:"source"`
	Service    string `yaml:"service" json:"service,omitempty"`
	SourcePath string `yaml:"source_path" json:"souce_path"`
}

// SetDefaults sets the default values.
func (d *AppMount) SetDefaults() {
	if d.Source == "" {
		d.Source = "local"
	}
}

// Validate checks for errors.
func (d AppMount) Validate(root *App) []error {
	o := make([]error, 0)
	if err := validateMustContainOne(
		[]string{"local", "service"},
		d.Source,
		"app.mounts[].source",
	); err != nil {
		o = append(o, err)
	}
	return o
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *AppMount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full app mount def
	data := make(map[string]string)
	e := unmarshal(&data)
	if e == nil {
		d.Source = data["source"]
		d.SourcePath = data["source_path"]
		d.Service = data["service"]
		return nil
	}
	// unmarshal string source path
	sourcePath := ""
	e = unmarshal(&sourcePath)
	if e != nil {
		return e
	}
	d.Source = "local"
	d.SourcePath = sourcePath
	mountStringStripPrefix := "shared:files"
	d.SourcePath = strings.TrimPrefix(d.SourcePath, mountStringStripPrefix)
	return nil
}
