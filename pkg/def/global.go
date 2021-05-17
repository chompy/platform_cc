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

// GlobalConfig defines global PCC configuration.
type GlobalConfig struct {
	Variables Variables         `yaml:"variables" json:"variables"`
	Flags     []string          `yaml:"flags" json:"flags"`
	Options   map[string]string `yaml:"options" json:"options"`
	Router    struct {
		PortHTTP  uint16 `yaml:"port_http" json:"port_http"`
		PortHTTPS uint16 `yaml:"port_https" json:"port_https"`
	} `yaml:"router" json:"router"`
}

// SetDefaults sets the default values.
func (d *GlobalConfig) SetDefaults() {
	if d.Router.PortHTTP == 0 {
		d.Router.PortHTTP = 80
	}
	if d.Router.PortHTTPS == 0 {
		d.Router.PortHTTPS = 443
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
								"global.variables.env."+k,
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
					"global.variables.env",
					"should be a map",
				))
				break
			}
		}
	}
	return o
}
