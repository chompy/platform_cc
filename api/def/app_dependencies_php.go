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

// AppDependenciesPhp defines php app dependencies.
type AppDependenciesPhp struct {
	Require      map[string]string               `yaml:"require" json:"require,omitempty"`
	Repositories []*AppDependenciesPhpRepository `yaml:"repositories" json:"repositories,omitempty"`
}

// SetDefaults sets the default values.
func (d *AppDependenciesPhp) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppDependenciesPhp) Validate(root *App) []error {
	o := make([]error, 0)
	for i := range d.Repositories {
		if e := d.Repositories[i].Validate(root); len(e) > 0 {
			o = append(o, e...)
		}
	}
	return o
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *AppDependenciesPhp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full extension
	data := make(map[string]interface{})
	e := unmarshal(&data)
	if e != nil {
		return e
	}
	d.Require = make(map[string]string)
	d.Repositories = make([]*AppDependenciesPhpRepository, 0)
	// dependencies as list of requirements with no repositories
	if data["require"] == nil {
		for k, v := range data {
			d.Require[k] = v.(string)
		}
		return nil
	}
	// includes repositories
	require := data["require"].(map[string]interface{})
	for k, v := range require {
		d.Require[k] = v.(string)
	}
	repos := data["repositories"].([]interface{})
	for _, v := range repos {
		d.Repositories = append(
			d.Repositories,
			&AppDependenciesPhpRepository{
				Type: v.(map[string]interface{})["type"].(string),
				URL:  v.(map[string]interface{})["url"].(string),
			},
		)
	}
	return nil
}
