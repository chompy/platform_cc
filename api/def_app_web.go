package api

// AppWebDef - defines how app is exposed to the web
type AppWebDef struct {
	Commands   AppWebCommandsDef             `yaml:"commands" json:"commands"`
	Upstream   AppWebUpstreamDef             `yaml:"upstream" json:"-"`
	Locations  map[string]*AppWebLocationDef `yaml:"locations" json:"locations"`
	MoveToRoot bool                          `json:"move_to_root"`
}

// SetDefaults - set default values
func (d *AppWebDef) SetDefaults() {
	d.Commands.SetDefaults()
	d.Upstream.SetDefaults()
	for i := range d.Locations {
		d.Locations[i].SetDefaults()
	}
	d.MoveToRoot = false
}

// Validate - validate AppWebDef
func (d AppWebDef) Validate() []error {
	o := make([]error, 0)
	if e := d.Commands.Validate(); e != nil {
		o = append(o, e...)
	}
	if e := d.Upstream.Validate(); e != nil {
		o = append(o, e...)
	}
	for _, l := range d.Locations {
		if e := l.Validate(); e != nil {
			o = append(o, e...)
		}
	}
	return o
}
