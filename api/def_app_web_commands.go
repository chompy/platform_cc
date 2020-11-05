package api

// AppWebCommandsDef - defines command to launch the app
type AppWebCommandsDef struct {
	Start string `yaml:"start" json:"start,omitempty"`
}

// SetDefaults - set default values
func (d *AppWebCommandsDef) SetDefaults() {
	return
}

// Validate - validate AppWebCommandsDef
func (d AppWebCommandsDef) Validate() []error {
	return []error{}
}
