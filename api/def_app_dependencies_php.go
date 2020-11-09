package api

// AppDependenciesPhpDef - define php app dependencies
type AppDependenciesPhpDef struct {
	Require      map[string]string                  `yaml:"require" json:"require,omitempty"`
	Repositories []*AppDependenciesPhpRepositoryDef `yaml:"repositories" json:"repositories,omitempty"`
}

// SetDefaults - set default values
func (d *AppDependenciesPhpDef) SetDefaults() {
	return
}

// Validate - validate AppDependenciesPhpDef
func (d AppDependenciesPhpDef) Validate() []error {
	o := make([]error, 0)
	for i := range d.Repositories {
		if e := d.Repositories[i].Validate(); len(e) > 0 {
			o = append(o, e...)
		}
	}
	return o
}

// UnmarshalYAML - implement Unmarshaler interface
func (d *AppDependenciesPhpDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full extension
	data := make(map[string]interface{})
	e := unmarshal(&data)
	if e != nil {
		return e
	}
	d.Require = make(map[string]string)
	d.Repositories = make([]*AppDependenciesPhpRepositoryDef, 0)
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
			&AppDependenciesPhpRepositoryDef{
				Type: v.(map[string]interface{})["type"].(string),
				URL:  v.(map[string]interface{})["url"].(string),
			},
		)
	}

	return nil
}
