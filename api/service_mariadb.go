package api

import (
	"strings"
)

// MariadbService - provides configuration for mariadb service
type MariadbService struct {
}

// getConfiguration - get configuration from yaml node
func (s MariadbService) getConfiguration(d *ServiceDef) (MariadbConfigurationDef, error) {
	o := MariadbConfigurationDef{}
	err := d.Configuration.Decode(&o)
	o.SetDefaults()
	return o, err
}

// CheckType - check if given type is used for this service
func (s MariadbService) CheckType(_type string) bool {
	return strings.HasPrefix(_type, "mariadb") || strings.HasPrefix(_type, "mysql")
}

// Validate - validate service definition
func (s MariadbService) Validate(d *ServiceDef) []error {
	c, err := s.getConfiguration(d)
	if err != nil {
		return []error{err}
	}
	return c.Validate()
}

// GetContainerConfig - get object with container configuration instructions
func (s MariadbService) GetContainerConfig(d *ServiceDef) ServiceContainerDef {
	config, _ := s.getConfiguration(d)
	return MariadbContainer{configuration: &config, def: d}
}
