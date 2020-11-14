package def

// AppWebLocation defines how the app should respond to a web request.
type AppWebLocation struct {
	Root             string                         `yaml:"root" json:"root"`
	Passthru         string                         `yaml:"passthru" json:"passthru"`
	Index            []string                       `yaml:"index" json:"index,omitempty"`
	Expires          string                         `yaml:"expires" json:"expires"`
	Scripts          Bool                           `yaml:"scripts" json:"scripts"`
	Allow            Bool                           `yaml:"allow" json:"allow"`
	Headers          map[string]string              `yaml:"headers" json:"headers,omitempty"`
	Rules            map[string]*AppWebLocation     `yaml:"rules" json:"rules,omitempty"`
	RequestBuffering AppWebLocationRequestBuffering `yaml:"request_buffering" json:"request_buffering"`
}

// SetDefaults sets the default values.
func (d *AppWebLocation) SetDefaults() {
	d.Scripts.DefaultValue = false
	if d.Passthru != "" && d.Passthru != "false" {
		d.Scripts.DefaultValue = true
	}
	if d.Expires == "" {
		d.Expires = "0"
	}
	d.Scripts.SetDefaults()
	d.Allow.DefaultValue = true
	d.Allow.SetDefaults()
	d.RequestBuffering.SetDefaults()
	for i := range d.Rules {
		d.Rules[i].SetDefaults()
	}
}

// Validate checks for errors.
func (d AppWebLocation) Validate(root *App) []error {
	o := make([]error, 0)
	// TODO validate expires
	// TODO validate headers?
	for _, r := range d.Rules {
		if e := r.Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
