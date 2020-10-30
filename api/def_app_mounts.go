package api

// AppMountDef - defines persistent mount volumes
type AppMountDef struct {
	Source     string `yaml:"source" json:"source"`
	Service    string `yaml:"service" json:"service,omitempty"`
	SourcePath string `yaml:"source_path" json:"souce_path"`
}

// SetDefaults - set default valuts
func (d *AppMountDef) SetDefaults() {
	if d.Source == "" {
		d.Source = "local"
	}
}

// Validate - validate AppMountDef
func (d AppMountDef) Validate() []error {
	o := make([]error, 0)
	if d.Source != "local" && d.Source != "service" {
		o = append(o, NewDefValidateError(
			"app.mounts[].source",
			"must be either local or service",
		))
	}
	return o
}

// UnmarshalYAML - implement Unmarshaler interface
func (d *AppMountDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full app mount def
	data := make(map[string]string)
	e := unmarshal(&data)
	if e == nil {
		d.Source = data["source"]
		d.SourcePath = data["source_path"]
		d.Service = data["service"]
		return nil
	}
	// unmarshal string source path
	sourcePath := ""
	e = unmarshal(&sourcePath)
	if e != nil {
		return e
	}
	d.Source = "local"
	d.SourcePath = sourcePath
	return nil
}
