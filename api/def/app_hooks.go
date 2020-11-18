package def

// AppHooks defines hook commands.
type AppHooks struct {
	Build      string `yaml:"build" json:"build"`
	Deploy     string `yaml:"deploy" json:"deploy"`
	PostDeploy string `yaml:"post_deploy" json:"post_deploy"`
}

// SetDefaults sets the default values.
func (d *AppHooks) SetDefaults() {
	return
}

// Validate checks for errors.
func (d AppHooks) Validate(root *App) []error {
	return []error{}
}
