package def

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ServiceConfiguration define service configuration.
type ServiceConfiguration map[string]interface{}

// UnmarshalYAML - parse yaml
func (d *ServiceConfiguration) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag != "!!map" {
		return fmt.Errorf("expected map in service configuration yaml")
	}
	*d = unmarshalYamlWithCustomTags(value).(map[string]interface{})
	return nil
}
