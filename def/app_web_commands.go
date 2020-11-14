package def

// AppWebCommands defines command(s) to launch the app.
type AppWebCommands struct {
	Start string `yaml:"start" json:"start,omitempty"`
}

// SetDefaults sets the default values.
func (d *AppWebCommands) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppWebCommands) Validate(root *App) []error {
	return []error{}
}
