package api

// AppHooksDef - defines hook commands
type AppHooksDef struct {
	Build      string `yaml:"build" json:"build"`
	Deploy     string `yaml:"deploy" json:"_deploy"`
	PostDeploy string `yaml:"post_deploy" json:"post_deploy"`
}

// SetDefaults - set default values
func (d *AppHooksDef) SetDefaults() {
	return
}

// Validate - validate AppHooksDef
func (d AppHooksDef) Validate() []error {
	return []error{}
}
