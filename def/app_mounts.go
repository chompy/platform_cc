package def

// AppMount defines persistent mount volumes
type AppMount struct {
	Source     string `yaml:"source" json:"source"`
	Service    string `yaml:"service" json:"service,omitempty"`
	SourcePath string `yaml:"source_path" json:"souce_path"`
}

// SetDefaults - set default valuts
func (d *AppMount) SetDefaults() {
	if d.Source == "" {
		d.Source = "local"
	}
}

// Validate checks for errors.
func (d AppMount) Validate(root *App) []error {
	o := make([]error, 0)
	if err := validateMustContainOne(
		[]string{"local", "service"},
		d.Source,
		"app.mounts[].source",
	); err != nil {
		o = append(o, err)
	}
	return o
}

// UnmarshalYAML implements Unmarshaler interface.
func (d *AppMount) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
