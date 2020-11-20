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

import "encoding/json"

// Bool defines a boolean that can contain a default value.
type Bool struct {
	value        bool
	isSet        bool
	DefaultValue bool
}

// SetDefaults sets the default value if not set via yaml.
func (d *Bool) SetDefaults() {
	if !d.isSet {
		d.value = d.DefaultValue
	}
}

// Get retrieves the current value.
func (d Bool) Get() bool {
	return d.value
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *Bool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&d.value)
	d.isSet = true
	return nil
}

// MarshalJSON implements json Marshaler interface.
func (d *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Get())
}
