package api

// ServicesMariadbConfigurationEndpointDef - define mariadb endpoint
type ServicesMariadbConfigurationEndpointDef struct {
	DefaultSchema string            `yaml:"default_schema"`
	Privileges    map[string]string `yaml:"privileges"`
}

// SetDefaults - set default values
func (d *ServicesMariadbConfigurationEndpointDef) SetDefaults() {
	return
}

// Validate - validate ServicesMariadbConfigurationEndpointDef
func (d ServicesMariadbConfigurationEndpointDef) Validate() []error {
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

func mariadbConfigurationEndpointFromInterface(i interface{}) *ServicesMariadbConfigurationEndpointDef {
	privFrom := i.(map[interface{}]interface{})["privileges"].(map[interface{}]interface{})
	privTo := make(map[string]string)
	for key, priv := range privFrom {
		privTo[key.(string)] = priv.(string)
	}
	return &ServicesMariadbConfigurationEndpointDef{
		DefaultSchema: i.(map[interface{}]interface{})["default_schema"].(string),
		Privileges:    privTo,
	}
}
