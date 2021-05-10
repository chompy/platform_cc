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

package platformsh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Environment defines a Platform.sh environment.
type Environment struct {
	Name         string `json:"name"`
	Title        string `json:"title"`
	IsMain       bool   `json:"is_main"`
	Status       string `json:"status"`
	MachineName  string `json:"machine_name"`
	EdgeHostname string `json:"edge_hostname"`
	Links        struct {
		SSHApp struct {
			HREF string `json:"href"`
		} `json:"pf:ssh:app"`
		SSH struct {
			HREF string `json:"href"`
		} `json:"ssh"`
		PublicURL struct {
			HREF string `json:"href"`
		} `json:"public-url"`
	} `json:"_links"`
}

// GetEnvironment returns environment matching given name.
func (p *Project) GetEnvironment(name string) *Environment {
	for i, e := range p.Environments {
		if e.Name == name || e.MachineName == name {
			return &p.Environments[i]
		}
	}
	return nil
}

// Variables returns list of variables for given platform.sh environment.
func (p *Project) Variables(env *Environment, service string) (map[string]string, error) {
	if env == nil {
		return nil, errors.WithStack(ErrInvalidEnvironment)
	}
	// fetch project variables
	projVarResp := make([]map[string]interface{}, 0)
	p.request(
		fmt.Sprintf("projects/%s/variables", p.ID),
		nil,
		&projVarResp,
	)
	// fetch environment variables
	envVarResp := make([]map[string]interface{}, 0)
	if err := p.request(
		fmt.Sprintf("projects/%s/environments/%s/variables", p.ID, env.Name),
		nil,
		&envVarResp,
	); err != nil {
		return nil, errors.WithStack(err)
	}
	// compile output
	out := make(map[string]string)
	for _, v := range append(projVarResp, envVarResp...) {
		if v["name"] == nil {
			continue
		}
		name := v["name"].(string)
		value := ""
		if v["value"] != nil {
			value = v["value"].(string)
		}
		// fetch var from environment if not found (i.e. sensitive)
		if value == "" && strings.HasPrefix(name, "env:") {
			envVarName := strings.TrimPrefix(name, "env:")
			var err error
			value, err = p.EnvironmentVariable(env, service, envVarName)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		out[name] = value
	}
	return out, nil
}

// EnvironmentVariable returns the value of the given environment variable.
func (p *Project) EnvironmentVariable(env *Environment, service string, name string) (string, error) {
	if env == nil {
		return "", errors.WithStack(ErrInvalidEnvironment)
	}
	out, err := p.SSHCommand(env, service, fmt.Sprintf("echo \"$%s\"", name))
	return string(out), errors.WithStack(err)
}

func (p *Project) decodeMapEnvVar(env *Environment, service string, name string) (map[string]interface{}, error) {
	// grab environment variable
	v, err := p.EnvironmentVariable(env, service, name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// decode base64
	b64DecOut, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// decode json
	out := make(map[string]interface{})
	if err := json.Unmarshal(b64DecOut, &out); err != nil {
		return nil, errors.WithStack(err)
	}
	return out, nil
}

// PlatformRelationships returns the value of then PLATFORM_RELATIONSHIPS environment variable.
func (p *Project) PlatformRelationships(env *Environment, service string) (map[string]interface{}, error) {
	return p.decodeMapEnvVar(env, service, "PLATFORM_RELATIONSHIPS")
}

// PlatformVariables returns the value of then PLATFORM_VARIABLES environment variable.
func (p *Project) PlatformVariables(env *Environment, service string) (map[string]interface{}, error) {
	return p.decodeMapEnvVar(env, service, "PLATFORM_VARIABLES")
}
