package api

import "encoding/json"

// BoolDef - boolean value
type BoolDef struct {
	value        bool
	isSet        bool
	DefaultValue bool
}

// SetDefaults - set default value if not set via yaml
func (d *BoolDef) SetDefaults() {
	if !d.isSet {
		d.value = d.DefaultValue
	}
}

// Validate - validate BoolDef
func (d BoolDef) Validate() []error {
	return []error{}
}

// Get - retrieve value
func (d BoolDef) Get() bool {
	return d.value
}

// UnmarshalYAML - implement Unmarshaler interface
func (d *BoolDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	unmarshal(&d.value)
	d.isSet = true
	return nil
}

// MarshalJSON - implement json Marshaler interface
func (d *BoolDef) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Get())
}
