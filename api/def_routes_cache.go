package api

// RouteCacheDef - defines caching for a route
type RouteCacheDef struct {
	Enabled    BoolDef  `yaml:"enabled" json:"enabled"`
	Headers    []string `yaml:"headers" json:"headers"`
	Cookies    []string `yaml:"cookies" json:"cookies"`
	DefaultTTL int      `yaml:"default_ttl" json:"default_ttl"`
}

// SetDefaults - set default values
func (d *RouteCacheDef) SetDefaults() {
	d.Enabled.DefaultValue = false
	d.Enabled.SetDefaults()
	if d.DefaultTTL == 0 {
		d.DefaultTTL = 300
	}
}

// Validate - validate RouteCacheDef
func (d RouteCacheDef) Validate() []error {
	return []error{}
}
