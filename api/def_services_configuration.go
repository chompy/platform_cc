package api

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ServiceConfigurationDef - define service configuration
type ServiceConfigurationDef map[string]interface{}

// UnmarshalYAML - parse yaml
func (d *ServiceConfigurationDef) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag != "!!map" {
		return fmt.Errorf("expected map in service configuration yaml")
	}
	*d = unmarshalYamlWithCustomTags(value).(map[string]interface{})
	return nil
}
