package api

// RoutesSsiDef - define route server side include
type RoutesSsiDef struct {
	Enabled BoolDef `yaml:"enabled"`
}

// SetDefaults - set default values
func (d *RoutesSsiDef) SetDefaults() {
	d.Enabled.DefaultValue = false
	d.Enabled.SetDefaults()
	return
}

// Validate - validate RoutesSsiDef
func (d RoutesSsiDef) Validate() []error {
	return []error{}
}
