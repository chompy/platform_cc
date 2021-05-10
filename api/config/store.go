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

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

const globalConfigPath = "config.json"

// Load loads the global configuration.
func Load() (def.GlobalConfig, error) {
	done := output.Duration("Load global configuration.")
	out := def.GlobalConfig{}
	// load config json
	raw, err := ioutil.ReadFile(pathTo(globalConfigPath))
	if err != nil {
		if os.IsNotExist(err) {
			out.Variables = make(def.Variables)
			out.Flags = make([]string, 0)
			out.Options = make(map[string]string)
			out.SetDefaults()
			done()
			return out, nil
		}
		return out, errors.WithStack(err)
	}
	// parse json
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, errors.WithStack(err)
	}
	out.SetDefaults()
	done()
	return out, nil
}

// Save saves given global config to file.
func Save(c def.GlobalConfig) error {
	done := output.Duration("Save global configuration.")
	// init config directory
	if err := initConfig(); err != nil {
		return errors.WithStack(err)
	}
	// encode json
	out, err := json.Marshal(c)
	if err != nil {
		return errors.WithStack(err)
	}
	// write to file
	if err := ioutil.WriteFile(pathTo(globalConfigPath), out, configPerm); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}
