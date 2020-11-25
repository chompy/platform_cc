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

// AppWorkersCommands defines command(s) to launch the app worker.
type AppWorkersCommands struct {
	Start string `yaml:"start" json:"start"`
}

// SetDefaults sets the default values.
func (d *AppWorkersCommands) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppWorkersCommands) Validate(root *App) []error {
	if d.Start == "" {
		return []error{
			NewValidateError("app.workers[].commands.start", "must not be empty"),
		}
	}
	return []error{}
}
