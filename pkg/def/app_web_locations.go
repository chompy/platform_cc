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

// AppWebLocation defines how the app should respond to a web request.
type AppWebLocation struct {
	Root             string                         `yaml:"root" json:"root"`
	Passthru         BoolString                     `yaml:"passthru" json:"passthru"`
	Index            []string                       `yaml:"index" json:"index,omitempty"`
	Expires          string                         `yaml:"expires" json:"expires"`
	Scripts          Bool                           `yaml:"scripts" json:"scripts"`
	Allow            Bool                           `yaml:"allow" json:"allow"`
	Headers          map[string]string              `yaml:"headers" json:"headers,omitempty"`
	Rules            map[string]*AppWebLocation     `yaml:"rules" json:"rules,omitempty"`
	RequestBuffering AppWebLocationRequestBuffering `yaml:"request_buffering" json:"request_buffering"`
}

// SetDefaults sets the default values.
func (d *AppWebLocation) SetDefaults() {
	d.Passthru.SetDefaults()
	d.Scripts.DefaultValue = false
	if d.Passthru.IsString() || d.Passthru.GetBool() {
		d.Scripts.DefaultValue = true
	}
	if d.Expires == "" {
		d.Expires = "0"
	}
	d.Scripts.SetDefaults()
	d.Allow.DefaultValue = true
	d.Allow.SetDefaults()
	d.RequestBuffering.SetDefaults()
	for i := range d.Rules {
		if !d.Rules[i].Passthru.IsSet() {
			d.Rules[i].Passthru = d.Passthru
		}
		d.Rules[i].SetDefaults()
	}
}

// Validate checks for errors.
func (d AppWebLocation) Validate(root *App) []error {
	o := make([]error, 0)
	// TODO validate expires
	// TODO validate headers?
	for _, r := range d.Rules {
		if e := r.Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
