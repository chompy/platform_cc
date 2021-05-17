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

// RouteCache defines caching for a route.
type RouteCache struct {
	Enabled    Bool     `yaml:"enabled" json:"enabled"`
	Headers    []string `yaml:"headers" json:"headers"`
	Cookies    []string `yaml:"cookies" json:"cookies"`
	DefaultTTL int      `yaml:"default_ttl" json:"default_ttl"`
}

// SetDefaults sets the default values.
func (d *RouteCache) SetDefaults() {
	d.Enabled.DefaultValue = false
	d.Enabled.SetDefaults()
	if d.DefaultTTL == 0 {
		d.DefaultTTL = 300
	}
}

// Validate checks for errors.
func (d RouteCache) Validate(root *Route) []error {
	return []error{}
}
