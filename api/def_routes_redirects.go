package api

// RouteRedirectsDef - define route redirects
type RouteRedirectsDef struct {
	Expires string                            `yaml:"expires"`
	Paths   map[string]*RouteRedirectsPathDef `yaml:"paths"`
}

// SetDefaults - set default values
func (d *RouteRedirectsDef) SetDefaults() {
	if d.Expires == "" {
		d.Expires = "-1"
	}
	for k := range d.Paths {
		d.Paths[k].SetDefaults()
	}
}

// Validate - validate RouteRedirectsDef
func (d RouteRedirectsDef) Validate() []error {
	o := make([]error, 0)
	for k := range d.Paths {
		if e := d.Paths[k].Validate(); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
