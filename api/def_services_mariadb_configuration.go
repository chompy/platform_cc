package api

// ServicesMariadbConfigurationDef - define configuration for mariadb service
type ServicesMariadbConfigurationDef struct {
	Schemas   []string                                            `yaml:"schemas"`
	Endpoints map[string]*ServicesMariadbConfigurationEndpointDef `yaml:"endpoints"`
}

// SetDefaults - set default values
func (d *ServicesMariadbConfigurationDef) SetDefaults() {
	for k := range d.Endpoints {
		d.Endpoints[k].SetDefaults()
	}
}

// Validate - validate ServicesMariadbConfigurationDef
func (d ServicesMariadbConfigurationDef) Validate() []error {
	o := make([]error, 0)
	for k := range d.Endpoints {
		if e := d.Endpoints[k].Validate(); e != nil {
			o = append(o, e...)
		}
	}
	return o
}

func mariadbConfigurationFromInterface(i interface{}) ServicesMariadbConfigurationDef {
	schemaFrom := i.(map[interface{}]interface{})["schemas"].([]interface{})
	schemaTo := make([]string, 0)
	for _, schema := range schemaFrom {
		schemaTo = append(schemaTo, schema.(string))
	}
	endpointsFrom := i.(map[interface{}]interface{})["endpoints"].(map[interface{}]interface{})
	endpointsTo := make(map[string]*ServicesMariadbConfigurationEndpointDef)
	for key, endpoint := range endpointsFrom {
		endpointsTo[key.(string)] = mariadbConfigurationEndpointFromInterface(endpoint)
	}
	return ServicesMariadbConfigurationDef{
		Schemas:   schemaTo,
		Endpoints: endpointsTo,
	}
}
