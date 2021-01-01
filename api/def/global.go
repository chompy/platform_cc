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

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gopkg.in/yaml.v3"
)

var globalConfigPaths = []string{
	"~/.config/platform_cc.yaml",
	"~/platform_cc.yaml",
}

// GlobalConfig defines global PCC configuration.
type GlobalConfig struct {
	Variables map[string]map[string]interface{} `yaml:"variables"`
	Router    struct {
		HTTP  uint16 `yaml:"http"`
		HTTPS uint16 `yaml:"https"`
	} `yaml:"router"`
}

// SetDefaults sets the default values.
func (d *GlobalConfig) SetDefaults() {
	if d.Router.HTTP == 0 {
		d.Router.HTTP = 80
	}
	if d.Router.HTTPS == 0 {
		d.Router.HTTPS = 443
	}
}

// Validate checks for errors.
func (d GlobalConfig) Validate() []error {
	o := make([]error, 0)
	return o
}

// ParseGlobalYaml parses the contents of a global configuration yaml file.
func ParseGlobalYaml(d []byte) (*GlobalConfig, error) {
	o := &GlobalConfig{
		Variables: make(map[string]map[string]interface{}),
	}
	err := yaml.Unmarshal(d, o)
	o.SetDefaults()
	return o, tracerr.Wrap(err)
}

// ParseGlobalYamlFile itterates list of possible global configuration yaml files and parse first one found.
func ParseGlobalYamlFile() (*GlobalConfig, error) {
	for _, gcp := range globalConfigPaths {
		d, err := ioutil.ReadFile(expandPath(gcp))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, tracerr.Wrap(err)
		}
		done := output.Duration(
			fmt.Sprintf("Parse global configuration at '%s.'", gcp),
		)
		o, err := ParseGlobalYaml(d)
		if err != nil {
			return nil, err
		}
		done()
		return o, nil
	}
	return &GlobalConfig{
		Variables: make(map[string]map[string]interface{}),
	}, nil

}
