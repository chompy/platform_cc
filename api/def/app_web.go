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

// AppWeb defines how app is exposed to the web.
type AppWeb struct {
	Commands   AppWebCommands             `yaml:"commands" json:"commands"`
	Upstream   AppWebUpstream             `yaml:"upstream" json:"upstream"`
	Locations  map[string]*AppWebLocation `yaml:"locations" json:"locations"`
	MoveToRoot bool                       `json:"move_to_root"`
}

// SetDefaults sets the default values.
func (d *AppWeb) SetDefaults() {
	d.Commands.SetDefaults()
	d.Upstream.SetDefaults()
	if d.Locations == nil || len(d.Locations) == 0 {
		d.Locations = map[string]*AppWebLocation{
			"/": &AppWebLocation{
				Passthru: BoolString{boolVal: true, isSet: true},
			},
		}
	}
	for i := range d.Locations {
		d.Locations[i].SetDefaults()
	}
	d.MoveToRoot = false
}

// Validate checks for errors.
func (d AppWeb) Validate(root *App) []error {
	o := make([]error, 0)
	if e := d.Commands.Validate(root); e != nil {
		o = append(o, e...)
	}
	if e := d.Upstream.Validate(root); e != nil {
		o = append(o, e...)
	}
	for _, l := range d.Locations {
		if e := l.Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
