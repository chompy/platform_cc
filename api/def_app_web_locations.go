package api

// AppWebLocationDef - defines how the app should respond to a web request
type AppWebLocationDef struct {
	Root             string                            `yaml:"root" json:"root"`
	Passthru         string                            `yaml:"passthru" json:"passthru"`
	Index            []string                          `yaml:"index" json:"index,omitempty"`
	Expires          string                            `yaml:"expires" json:"expires"`
	Scripts          BoolDef                           `yaml:"scripts" json:"scripts"`
	Allow            BoolDef                           `yaml:"allow" json:"allow"`
	Headers          map[string]string                 `yaml:"headers" json:"headers,omitempty"`
	Rules            map[string]*AppWebLocationDef     `yaml:"rules" json:"rules,omitempty"`
	RequestBuffering AppWebLocationRequestBufferingDef `yaml:"request_buffering" json:"request_buffering"`
}

// SetDefaults - set default values
func (d *AppWebLocationDef) SetDefaults() {
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

// Validate - validate AppWebLocationDef
func (d AppWebLocationDef) Validate() []error {
	o := make([]error, 0)
	// TODO validate expires
	// TODO validate headers?
	for _, r := range d.Rules {
		if e := r.Validate(); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
