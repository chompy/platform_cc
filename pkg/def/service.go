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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
)

// Service defines a service.
type Service struct {
	Name          string
	Type          string               `yaml:"type" json:"type"`
	Disk          int                  `yaml:"disk" json:"disk"`
	Configuration ServiceConfiguration `yaml:"configuration" json:"configuration,omitempty"`
	Relationships map[string]string    `yaml:"relationships" json:"relationships,omitempty"`
	Disable       bool                 `yaml:"_disable"`
}

// SetDefaults sets the default values.
func (d *Service) SetDefaults() {
	if d.Configuration == nil {
		d.Configuration = make(map[string]interface{})
	}
	if d.Configuration["application_size"] == nil {
		d.Configuration["application_size"] = 0
	}
}

// Validate checks for errors.
func (d Service) Validate() []error {
	o := make([]error, 0)
	if d.Type == "" {
		o = append(o, NewValidateError(
			fmt.Sprintf("services.%s.type", d.Name),
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

// ParseServiceYamls parses multiple services.yaml contents and merges them in to one.
func ParseServiceYamls(d [][]byte) ([]Service, error) {
	o := make(map[string]*Service)
	if err := mergeYamls(d, o); err != nil {
		return []Service{}, errors.WithStack(err)
	}
	// set defaults + transfer to new slice
	oo := make([]Service, 0)
	for k := range o {
		if o[k].Disable {
			continue
		}
		o[k].SetDefaults()
		o[k].Name = k
		oo = append(oo, *o[k])
	}
	return oo, nil
}

// ParseServiceYamlFiles parses multiple services.yaml files and merges them in to one.
func ParseServiceYamlFiles(fileList []string) ([]Service, error) {
	done := output.Duration(
		fmt.Sprintf("Parse service at '%s.'", strings.Join(fileList, ", ")),
	)
	byteList := make([][]byte, 0)
	for _, f := range fileList {
		projectPlatformDir = filepath.Dir(f)
		d, err := ioutil.ReadFile(f)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, errors.WithStack(err)
		}
		byteList = append(byteList, d)
	}
	a, err := ParseServiceYamls(byteList)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	done()
	return a, nil
}
