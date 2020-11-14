package def

import "encoding/json"

// Bool defines a boolean that can contain a default value.
type Bool struct {
	value        bool
	isSet        bool
	DefaultValue bool
}

// SetDefaults - set default value if not set via yaml
func (d *Bool) SetDefaults() {
	if !d.isSet {
		d.value = d.DefaultValue
	}
}

// Get - retrieve value
func (d Bool) Get() bool {
	return d.value
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *Bool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&d.value)
	d.isSet = true
	return nil
}

// MarshalJSON - implement json Marshaler interface
func (d *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Get())
}
