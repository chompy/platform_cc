package def

// AppRuntime defines runtime configuration.
type AppRuntime struct {
	RequestTerminateTimeout int                    `yaml:"request_terminate_timeout" json:"-"`
	Extensions              []*AppRuntimeExtension `yaml:"extensions" json:"extensions"`
	DisabledExtensions      []string               `yaml:"disabled_extensions" json:"disabled_extensions,omitempty"`
}

// SetDefaults sets the default values.
func (d *AppRuntime) SetDefaults() {
	if d.RequestTerminateTimeout <= 0 {
		d.RequestTerminateTimeout = 300
	}
	if d.Extensions == nil || len(d.Extensions) == 0 {
		d.Extensions = make([]*AppRuntimeExtension, 0)
	}
	for i := range d.Extensions {
		d.Extensions[i].SetDefaults()
	}
}

// Validate checks for errors.
func (d AppRuntime) Validate(root *App) []error {
	o := make([]error, 0)
	for i := range d.Extensions {
		if e := d.Extensions[i].Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
