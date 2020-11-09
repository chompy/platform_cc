package api

// AppDependenciesDef - define dependencies of the application
type AppDependenciesDef struct {
	PHP     AppDependenciesPhpDef `yaml:"php" json:"php"`
	NodeJS  map[string]string     `yaml:"nodejs" json:"nodejs"`
	Python2 map[string]string     `yaml:"python2" json:"python2"`
	Python3 map[string]string     `yaml:"python3" json:"python3"`
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
