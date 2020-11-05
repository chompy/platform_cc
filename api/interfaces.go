package api

// Def - define definition interface
type Def interface {
	Validate() []error
}

// ServiceResolver - resolve service configuration
type ServiceResolver interface {
	CheckType(string) bool
	Validate(*ServiceDef) []error
	GetContainerConfig(*ServiceDef) ServiceContainerDef
}

// ServiceContainerDef - define service container
type ServiceContainerDef interface {
	GetImage() string
	GetVolumes() []string
	GetStartCommand() []string
	GetPostStartCommand() []string
	GetRelationship() []map[string]interface{}
}
