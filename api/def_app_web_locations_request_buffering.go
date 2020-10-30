package api

// AppWebLocationRequestBufferingDef - defines request buffering config
type AppWebLocationRequestBufferingDef struct {
	Enabled        BoolDef `yaml:"enabled" json:"enabled"`
	MaxRequestSize string  `yaml:"max_request_size" json:"max_request_size"`
}

// SetDefaults - set default values
func (d *AppWebLocationRequestBufferingDef) SetDefaults() {
	d.Enabled.DefaultValue = true
	d.Enabled.SetDefaults()
	if d.MaxRequestSize == "" {
		d.MaxRequestSize = "250m"
	}
}

// Validate - validate AppWebLocationRequestBufferingDef
func (d AppWebLocationRequestBufferingDef) Validate() []error {
	// TODO validate max request size
	return []error{}
}
