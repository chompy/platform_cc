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
	"os/user"
	"strings"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gopkg.in/yaml.v3"
)

var globalConfigPaths = []string{
	"~/.pcc/config.yaml",
	"~/.config/platform_cc.yaml",
	"~/platform_cc.yaml",
}

const defaultSSHKeyPath = "~/.ssh/pccid"

// GlobalConfig defines global PCC configuration.
type GlobalConfig struct {
	Variables Variables         `yaml:"variables"`
	Flags     []string          `yaml:"flags"`
	Options   map[string]string `yaml:"options"`
	Router    struct {
		HTTP  uint16 `yaml:"http"`
		HTTPS uint16 `yaml:"https"`
	} `yaml:"router"`
	SSH struct {
		KeyPath string `yaml:"key_path"`
		Key     string `yaml:"key"`
	} `yaml:"ssh"`
	PlatformSH struct {
		APIToken string `yaml:"api_token"`
	} `yaml:"platform_sh"`
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
	if d.Variables["env"] != nil {
		switch d.Variables["env"].(type) {
		case map[string]interface{}:
			{
				for k, v := range d.Variables["env"].(map[string]interface{}) {
					switch v.(type) {
					case string, int, float32, float64, bool:
						{
							break
						}
					default:
						{
							o = append(o, NewValidateError(
								"app.variables.env."+k,
								"should be a scalar",
							))
							break
						}
					}
				}
				break
			}
		default:
			{
				o = append(o, NewValidateError(
					"app.variables.env",
					"should be a map",
				))
				break
			}
		}
	}
	return o
}

// GetSSHKey returns SSH key.
func (d GlobalConfig) GetSSHKey() string {
	keyPath := defaultSSHKeyPath
	if d.SSH.Key != "" {
		return d.SSH.Key
	} else if d.SSH.KeyPath != "" {
		keyPath = d.SSH.KeyPath
	}
	// get user home dir
	homeDir := "~"
	user, err := user.Current()
	if err != nil {
		output.LogError(err)
	} else {
		homeDir = user.HomeDir
	}
	// fix path
	keyPath = strings.ReplaceAll(keyPath, "~", strings.TrimRight(homeDir, string(os.PathSeparator)))
	// read key
	sshKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		output.LogError(err)
	}
	return string(sshKey)
}

// ParseGlobalYaml parses the contents of a global configuration yaml file.
func ParseGlobalYaml(d []byte) (*GlobalConfig, error) {
	o := &GlobalConfig{
		Variables: make(Variables),
	}
	err := yaml.Unmarshal(d, o)
	o.SetDefaults()
	return o, tracerr.Wrap(err)
}

// ParseGlobalYamlFile itterates list of possible global configuration yaml files and parse first one found.
func ParseGlobalYamlFile() (*GlobalConfig, error) {
	o := &GlobalConfig{
		Variables: make(Variables),
	}
	for _, gcp := range globalConfigPaths {
		d, err := ioutil.ReadFile(expandPath(gcp))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return o, tracerr.Wrap(err)
		}
		done := output.Duration(
			fmt.Sprintf("Parse global configuration at '%s.'", gcp),
		)
		o, err := ParseGlobalYaml(d)
		if err != nil {
			return o, err
		}
		done()
		return o, nil
	}
	o.SetDefaults()
	return o, nil
}
