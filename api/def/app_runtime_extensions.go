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

// AppRuntimeExtension defines an extension (PHP).
type AppRuntimeExtension struct {
	Name          string            `yaml:"name" json:"name"`
	Configuration map[string]string `yaml:"configuration" json:"configuration,omitempty"`
}

// SetDefaults sets the default values.
func (d *AppRuntimeExtension) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppRuntimeExtension) Validate(root *App) []error {
	return []error{}
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *AppRuntimeExtension) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full extension
	data := make(map[string]interface{})
	e := unmarshal(&data)
	if e == nil {
		d.Name = data["name"].(string)
		d.Configuration = make(map[string]string)
		conf := data["configuration"].(map[string]interface{})
		for k, v := range conf {
			d.Configuration[k] = v.(string)
		}
		return nil
	}
	// unmarshal string extension name
	extName := ""
	e = unmarshal(&extName)
	if e != nil {
		return e
	}
	d.Name = extName
	d.Configuration = make(map[string]string)
	return nil
}
