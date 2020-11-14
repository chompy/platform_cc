package def

// AppDependencies defines dependencies of the application.
type AppDependencies struct {
	PHP     AppDependenciesPhp `yaml:"php" json:"php"`
	NodeJS  map[string]string  `yaml:"nodejs" json:"nodejs"`
	Python2 map[string]string  `yaml:"python2" json:"python2"`
	Python3 map[string]string  `yaml:"python3" json:"python3"`
}

// SetDefaults sets the default values.
func (d *AppDependencies) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppDependencies) Validate(root *App) []error {
	o := make([]error, 0)
	if e := d.PHP.Validate(root); len(e) > 0 {
		o = append(o, e...)
	}
	return o
}
