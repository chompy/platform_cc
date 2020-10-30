package api

// AppDependenciesDef - define dependencies of the application
type AppDependenciesDef struct {
	PHP AppDependenciesPhpDef `yaml:"php"`
}

// SetDefaults - set default values
func (d *AppDependenciesDef) SetDefaults() {
	return
}

// Validate - validate AppDependenciesDef
func (d AppDependenciesDef) Validate() []error {
	o := make([]error, 0)
	if e := d.PHP.Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	return o
}
