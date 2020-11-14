package def

// RouteCache defines caching for a route.
type RouteCache struct {
	Enabled    Bool     `yaml:"enabled" json:"enabled"`
	Headers    []string `yaml:"headers" json:"headers"`
	Cookies    []string `yaml:"cookies" json:"cookies"`
	DefaultTTL int      `yaml:"default_ttl" json:"default_ttl"`
}

// SetDefaults sets the default values.
func (d *RouteCache) SetDefaults() {
	d.Enabled.DefaultValue = false
	d.Enabled.SetDefaults()
	if d.DefaultTTL == 0 {
		d.DefaultTTL = 300
	}
}

// Validate checks for errors.
func (d RouteCache) Validate(root *Route) []error {
	return []error{}
}
