package def

// RoutesSsi define route server side include.
type RoutesSsi struct {
	Enabled Bool `yaml:"enabled" json:"enabled"`
}

// SetDefaults sets the default values.
func (d *RoutesSsi) SetDefaults() {
	d.Enabled.DefaultValue = false
	d.Enabled.SetDefaults()
	return
}

// Validate checks for errors.
func (d RoutesSsi) Validate(root *Route) []error {
	return []error{}
}
