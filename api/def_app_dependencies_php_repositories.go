package api

// AppDependenciesPhpRepositoryDef - define php app dependency repository
type AppDependenciesPhpRepositoryDef struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

// SetDefaults - set default values
func (d *AppDependenciesPhpRepositoryDef) SetDefaults() {
	if d.Type == "" {
		d.Type = "vcs"
	}
}

// Validate - validate AppDependenciesPhpRepositoryDef
func (d AppDependenciesPhpRepositoryDef) Validate() []error {
	return []error{}
}
