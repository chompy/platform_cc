package def

// AppDependenciesPhpRepository defines php app dependency repository.
type AppDependenciesPhpRepository struct {
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}

// SetDefaults sets the default values.
func (d *AppDependenciesPhpRepository) SetDefaults() {
	if d.Type == "" {
		d.Type = "vcs"
	}
}

// Validate checks for errors.
func (d AppDependenciesPhpRepository) Validate(root *App) []error {
	return []error{}
}
