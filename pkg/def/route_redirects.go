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

// RouteRedirects define route redirects.
type RouteRedirects struct {
	Expires string                         `yaml:"expires" json:"expires"`
	Paths   map[string]*RouteRedirectsPath `yaml:"paths" json:"paths"`
}

// SetDefaults sets the default values.
func (d *RouteRedirects) SetDefaults() {
	if d.Expires == "" {
		d.Expires = "-1"
	}
	for k := range d.Paths {
		d.Paths[k].SetDefaults()
	}
}

// Validate checks for errors.
func (d RouteRedirects) Validate(root *Route) []error {
	o := make([]error, 0)
	// TODO validate expires
	for k := range d.Paths {
		if e := d.Paths[k].Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
