package api

// MariadbConfigurationDef - define configuration for mariadb service
type MariadbConfigurationDef struct {
	Schemas   []string                                    `yaml:"schemas"`
	Endpoints map[string]*MariadbConfigurationEndpointDef `yaml:"endpoints"`
}

// SetDefaults - set default values
func (d *MariadbConfigurationDef) SetDefaults() {
	for k := range d.Endpoints {
		d.Endpoints[k].SetDefaults()
	}
	if len(d.Schemas) == 0 {
		d.Schemas = []string{"mysql"}
	}
	if len(d.Endpoints) == 0 {
		d.Endpoints["mysql"] = &MariadbConfigurationEndpointDef{
			DefaultSchema: d.Schemas[0],
			Privileges: map[string]string{
				d.Schemas[0]: "admin",
			},
		}
	}
}

// Validate - validate MariadbConfigurationDef
func (d MariadbConfigurationDef) Validate() []error {
	o := make([]error, 0)
	for k := range d.Endpoints {
		if e := d.Endpoints[k].Validate(); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
