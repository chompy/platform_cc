package api

// RouteRedirectsPathDef - defines a route redirect path
type RouteRedirectsPathDef struct {
	To           string  `yaml:"to" json:"to"`
	Regexp       BoolDef `yaml:"regexp" json:"regexp"`
	Prefix       BoolDef `yaml:"prefix" json:"prefix"`
	AppendSuffix BoolDef `yaml:"append_suffix" json:"append_suffix"`
	Code         int     `yaml:"code" json:"code"`
	Expires      string  `yaml:"expires" json:"expires"`
}

// SetDefaults - set default values
func (d *RouteRedirectsPathDef) SetDefaults() {
	d.Regexp.DefaultValue = false
	d.Regexp.SetDefaults()
	d.Prefix.DefaultValue = true
	d.Prefix.SetDefaults()
	d.AppendSuffix.DefaultValue = true
	d.AppendSuffix.SetDefaults()
	if d.Code == 0 {
		d.Code = 302
	}
}

// Validate - validate RouteRedirectsPathDef
func (d RouteRedirectsPathDef) Validate() []error {
	o := make([]error, 0)
	if d.Code != 301 && d.Code != 302 && d.Code != 307 && d.Code != 308 {
		o = append(o, NewDefValidateError(
			"routes[].redirects.paths[].code",
			"invalid status code, valid codes are 301, 302, 307, and 308",
		))
	}
	return o
}
