package api

// AppBuildDef - defines what happens when building the app
type AppBuildDef struct {
	Flavor string `yaml:"flavor"`
}

// SetDefaults - set default values
func (d *AppBuildDef) SetDefaults() {
	if d.Flavor == "" {
		d.Flavor = "none"
	}
}

// Validate - validate AppBuildDef
func (d AppBuildDef) Validate() []error {
	return []error{}
}
