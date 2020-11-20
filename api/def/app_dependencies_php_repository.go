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

// AppDependenciesPhpRepository defines php app dependency repository.
type AppDependenciesPhpRepository struct {
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}

// SetDefaults sets the default values.
func (d *AppDependenciesPhpRepository) SetDefaults() {
	if d.Type == "" {
		d.Type = "vcs"
	}
}

// Validate checks for errors.
func (d AppDependenciesPhpRepository) Validate(root *App) []error {
	return []error{}
}
