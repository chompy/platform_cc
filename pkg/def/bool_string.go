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

// BoolString defines a type that can either be a boolean or a string.
type BoolString struct {
	stringVal string
	boolVal   bool
	isSet     bool
}

// SetDefaults sets the default value if not set via yaml.
func (d *BoolString) SetDefaults() {
	if !d.isSet {
		d.boolVal = false
		d.stringVal = ""
	}
}

// IsString returns true if current value is string.
func (d BoolString) IsString() bool {
	return d.stringVal != ""
}

// GetString retrieves the current string value.
func (d BoolString) GetString() string {
	return d.stringVal
}

// GetBool retrieves the current bool value.
func (d BoolString) GetBool() bool {
	return d.boolVal
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *BoolString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val interface{}
	unmarshal(&val)
	switch val := val.(type) {
	case bool:
		{
			d.boolVal = val
			d.stringVal = ""
			break
		}
	case string:
		{
			d.boolVal = false
			d.stringVal = val
			break
		}
	}
	d.isSet = true
	return nil
}

// MarshalJSON implements json Marshaler interface.
func (d *BoolString) MarshalJSON() ([]byte, error) {
	if d.IsString() {
		return json.Marshal(d.GetString())
	}
	return json.Marshal(d.GetBool())
}
