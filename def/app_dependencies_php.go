package def

// AppDependenciesPhp defines php app dependencies.
type AppDependenciesPhp struct {
	Require      map[string]string               `yaml:"require" json:"require,omitempty"`
	Repositories []*AppDependenciesPhpRepository `yaml:"repositories" json:"repositories,omitempty"`
}

// SetDefaults sets the default values.
func (d *AppDependenciesPhp) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppDependenciesPhp) Validate(root *App) []error {
	o := make([]error, 0)
	for i := range d.Repositories {
		if e := d.Repositories[i].Validate(root); len(e) > 0 {
			o = append(o, e...)
		}
	}
	return o
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *AppDependenciesPhp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full extension
	data := make(map[string]interface{})
	e := unmarshal(&data)
	if e != nil {
		return e
	}
	d.Require = make(map[string]string)
	d.Repositories = make([]*AppDependenciesPhpRepository, 0)
	// dependencies as list of requirements with no repositories
	if data["require"] == nil {
		for k, v := range data {
			d.Require[k] = v.(string)
		}
		return nil
	}
	// includes repositories
	require := data["require"].(map[string]interface{})
	for k, v := range require {
		d.Require[k] = v.(string)
	}
	repos := data["repositories"].([]interface{})
	for _, v := range repos {
		d.Repositories = append(
			d.Repositories,
			&AppDependenciesPhpRepository{
				Type: v.(map[string]interface{})["type"].(string),
				URL:  v.(map[string]interface{})["url"].(string),
			},
		)
	}
	return nil
}
