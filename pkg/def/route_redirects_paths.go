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

// RouteRedirectsPath defines a route redirect path.
type RouteRedirectsPath struct {
	To           string `yaml:"to" json:"to"`
	Regexp       Bool   `yaml:"regexp" json:"regexp"`
	Prefix       Bool   `yaml:"prefix" json:"prefix"`
	AppendSuffix Bool   `yaml:"append_suffix" json:"append_suffix"`
	Code         int    `yaml:"code" json:"code"`
	Expires      string `yaml:"expires" json:"expires"`
}

// SetDefaults sets the default values.
func (d *RouteRedirectsPath) SetDefaults() {
	d.Regexp.DefaultValue = false
	d.Regexp.SetDefaults()
	d.Prefix.DefaultValue = true
	d.Prefix.SetDefaults()
	d.AppendSuffix.DefaultValue = true
	d.AppendSuffix.SetDefaults()
	if d.Code == 0 {
		d.Code = 302
	}
}

// Validate checks for errors.
func (d RouteRedirectsPath) Validate(root *Route) []error {
	o := make([]error, 0)
	if err := validateMustContainOneInt(
		[]int{301, 302, 307, 308},
		d.Code,
		"routes[].redirects.paths[].code",
	); err != nil {
		o = append(o, err)
	}
	return o
}
