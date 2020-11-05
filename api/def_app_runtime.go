package api

// AppRuntimeDef - defines runtime configuration
type AppRuntimeDef struct {
	RequestTerminateTimeout int                       `yaml:"request_terminate_timeout" json:"-"`
	Extensions              []*AppRuntimeExtensionDef `yaml:"extensions" json:"extensions"`
	DisabledExtensions      []string                  `yaml:"disabled_extensions" json:"disabled_extensions,omitempty"`
}

// SetDefaults - set default values
func (d *AppRuntimeDef) SetDefaults() {
	if d.RequestTerminateTimeout <= 0 {
		d.RequestTerminateTimeout = 300
	}
	if d.Extensions == nil || len(d.Extensions) == 0 {
		d.Extensions = make([]*AppRuntimeExtensionDef, 0)
	}
	for i := range d.Extensions {
		d.Extensions[i].SetDefaults()
	}
}

// Validate - validate AppRuntimeDef
func (d AppRuntimeDef) Validate() []error {
	o := make([]error, 0)
	for i := range d.Extensions {
		if e := d.Extensions[i].Validate(); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
