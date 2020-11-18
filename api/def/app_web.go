package def

// AppWeb defines how app is exposed to the web.
type AppWeb struct {
	Commands   AppWebCommands             `yaml:"commands" json:"commands"`
	Upstream   AppWebUpstream             `yaml:"upstream" json:"-"`
	Locations  map[string]*AppWebLocation `yaml:"locations" json:"locations"`
	MoveToRoot bool                       `json:"move_to_root"`
}

// SetDefaults sets the default values.
func (d *AppWeb) SetDefaults() {
	d.Commands.SetDefaults()
	d.Upstream.SetDefaults()
	for i := range d.Locations {
		d.Locations[i].SetDefaults()
	}
	d.MoveToRoot = false
}

// Validate checks for errors.
func (d AppWeb) Validate(root *App) []error {
	o := make([]error, 0)
	if e := d.Commands.Validate(root); e != nil {
		o = append(o, e...)
	}
	if e := d.Upstream.Validate(root); e != nil {
		o = append(o, e...)
	}
	for _, l := range d.Locations {
		if e := l.Validate(root); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
