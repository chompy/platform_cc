package api

// MariadbConfigurationEndpointDef - define mariadb endpoint
type MariadbConfigurationEndpointDef struct {
	DefaultSchema string            `yaml:"default_schema"`
	Privileges    map[string]string `yaml:"privileges"`
}

// SetDefaults - set default values
func (d *MariadbConfigurationEndpointDef) SetDefaults() {
	return
}

// Validate - validate MariadbConfigurationEndpointDef
func (d MariadbConfigurationEndpointDef) Validate() []error {
	o := make([]error, 0)
	for _, v := range d.Privileges {
		if v != "ro" && v != "rw" && v != "admin" {
			o = append(o, NewDefValidateError(
				"services[].configuration[mariadb].endpoints[].privileges",
				"must be ro, rw, or admin",
			))
		}
	}
	return o
}
