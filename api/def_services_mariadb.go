package api

// ServiceMariadbDef - defines mariadb service
type ServiceMariadbDef struct {
	*ServiceDef   `yaml:",inline"`
	Configuration ServicesMariadbConfigurationDef `yaml:"configuration"`
}

// SetDefaults - set default values
func (d *ServiceMariadbDef) SetDefaults() {
	d.Configuration.SetDefaults()
}

// Validate - validate
func (d ServiceMariadbDef) Validate() []error {
	o := make([]error, 0)
	if e := d.Configuration.Validate(); e != nil {
		o = append(o, e...)
	}
	return o
}
