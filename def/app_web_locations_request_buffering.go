package def

// AppWebLocationRequestBuffering defines request buffering config.
type AppWebLocationRequestBuffering struct {
	Enabled        Bool   `yaml:"enabled" json:"enabled"`
	MaxRequestSize string `yaml:"max_request_size" json:"max_request_size"`
}

// SetDefaults sets the default values.
func (d *AppWebLocationRequestBuffering) SetDefaults() {
	d.Enabled.DefaultValue = true
	d.Enabled.SetDefaults()
	if d.MaxRequestSize == "" {
		d.MaxRequestSize = "250m"
	}
}

// Validate checks for errors.
func (d AppWebLocationRequestBuffering) Validate(root *App) []error {
	// TODO validate max request size
	return []error{}
}
