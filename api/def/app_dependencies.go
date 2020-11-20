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

// AppDependencies defines dependencies of the application.
type AppDependencies struct {
	PHP     AppDependenciesPhp `yaml:"php" json:"php"`
	NodeJS  map[string]string  `yaml:"nodejs" json:"nodejs"`
	Python2 map[string]string  `yaml:"python2" json:"python2"`
	Python3 map[string]string  `yaml:"python3" json:"python3"`
}

// SetDefaults sets the default values.
func (d *AppDependencies) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppDependencies) Validate(root *App) []error {
	o := make([]error, 0)
	if e := d.PHP.Validate(root); len(e) > 0 {
		o = append(o, e...)
	}
	return o
}
