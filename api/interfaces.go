package api

// Def - define definition interface
type Def interface {
	Validate() []error
}

// ServiceConfig - service configuration
type ServiceConfig interface {
	Check(*ServiceDef) bool
	Validate(*ServiceDef) []error
	GetSetupCommand(*ServiceDef) ([]string, error)
	GetRelationship(*ServiceDef) ([]map[string]interface{}, error)
}
