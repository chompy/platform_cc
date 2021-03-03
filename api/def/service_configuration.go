/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package def

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ServiceConfiguration define service configuration.
type ServiceConfiguration map[string]interface{}

// IsAuthenticationEnabled returns true if 'authentication.enabled' is true.
func (d ServiceConfiguration) IsAuthenticationEnabled() bool {
	switch d["authentication"].(type) {
	case map[string]interface{}:
		{
			enabledConf := d["authentication"].(map[string]interface{})["enabled"]
			switch enabledConf.(type) {
			case bool:
				{
					return enabledConf.(bool)
				}
			case int:
				{
					return enabledConf.(int) != 0
				}
			case string:
				{
					enabledStr := enabledConf.(string)
					return enabledStr != "" && enabledStr != "0" && enabledStr != "no" && enabledStr != "false" && enabledStr != "off"
				}
			}
			return false
		}
	}
	return false
}

// UnmarshalYAML - parse yaml
func (d *ServiceConfiguration) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag != "!!map" {
		return fmt.Errorf("expected map in service configuration yaml")
	}
	*d = unmarshalYamlWithCustomTags(value).(map[string]interface{})
	return nil
}
