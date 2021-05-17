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

// AppWebLocationRequestBuffering defines request buffering config.
type AppWebLocationRequestBuffering struct {
	Enabled        Bool   `yaml:"enabled" json:"enabled"`
	MaxRequestSize string `yaml:"max_request_size" json:"max_request_size"`
}

// SetDefaults sets the default values.
func (d *AppWebLocationRequestBuffering) SetDefaults() {
	d.Enabled.DefaultValue = true
	d.Enabled.SetDefaults()
	if d.MaxRequestSize == "" {
		d.MaxRequestSize = "250m"
	}
}

// Validate checks for errors.
func (d AppWebLocationRequestBuffering) Validate(root *App) []error {
	// TODO validate max request size
	return []error{}
}
