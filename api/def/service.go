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
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Service defines a service.
type Service struct {
	Name          string
	Type          string               `yaml:"type" json:"type"`
	Disk          int                  `yaml:"disk" json:"disk"`
	Configuration ServiceConfiguration `yaml:"configuration" json:"configuration,omitempty"`
	Relationships map[string]string    `yaml:"relationships" json:"relationships,omitempty"`
}

// SetDefaults sets the default values.
func (d *Service) SetDefaults() {
	if d.Configuration == nil {
		d.Configuration = make(map[string]interface{})
	}
	if d.Configuration["application_size"] == nil {
		d.Configuration["application_size"] = 0
	}
	return
}

// Validate checks for errors.
func (d Service) Validate() []error {
	o := make([]error, 0)
	if d.Type == "" {
		o = append(o, NewValidateError(
			"services[].type",
			"must be defined",
		))
	}
	return o
}

// GetTypeName gets the service type.
func (d Service) GetTypeName() string {
	return strings.Split(d.Type, ":")[0]
}

// GetEmptyRelationship retursn an empty relationship.
func (d Service) GetEmptyRelationship() map[string]interface{} {
	return map[string]interface{}{
		"host":        "",
		"hostname":    "",
		"ip":          "",
		"port":        80,
		"path":        "",
		"scheme":      d.GetTypeName(),
		"fragment":    nil,
		"rel":         "",
		"host_mapped": false,
		"public":      false,
		"type":        d.Type,
		"service":     d.Name,
	}
}

// ParseServicesYaml parses the contents of services.yaml.
func ParseServicesYaml(d []byte) ([]Service, error) {
	o := make(map[string]*Service)
	e := yaml.Unmarshal(d, &o)
	oo := make([]Service, 0)
	for k := range o {
		o[k].SetDefaults()
		o[k].Name = k
		oo = append(oo, *o[k])
	}
	return oo, e
}

// ParseServicesYamlFile - open services yaml file and parse it
func ParseServicesYamlFile(f string) ([]Service, error) {
	log.Printf("Parse services at '%s.'", f)
	projectPlatformDir = filepath.Dir(f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return []Service{}, e
	}
	return ParseServicesYaml(d)
}
