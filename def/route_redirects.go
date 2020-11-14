package def

// RouteRedirects define route redirects.
type RouteRedirects struct {
	Expires string                         `yaml:"expires" json:"expires"`
	Paths   map[string]*RouteRedirectsPath `yaml:"paths" json:"paths"`
}

// SetDefaults sets the default values.
func (d *RouteRedirects) SetDefaults() {
	if d.Expires == "" {
		d.Expires = "-1"
	}
	for k := range d.Paths {
		d.Paths[k].SetDefaults()
	}
}

// Validate checks for errors.
func (d RouteRedirects) Validate(root *Route) []error {
	o := make([]error, 0)
	// TODO validate expires
	for k := range d.Paths {
		if e := d.Paths[k].Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
