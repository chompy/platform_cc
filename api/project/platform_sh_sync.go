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

package project

import (
	"fmt"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// PlatformSHSyncVariables syncs the given platform.sh environment's variables to the local project.
func (p *Project) PlatformSHSyncVariables(envName string) error {
	if p.PlatformSH == nil || p.PlatformSH.ID == "" {
		return tracerr.Errorf("platform.sh project not found")
	}
	if len(p.Apps) < 1 {
		return tracerr.Errorf("project should have at least one application")
	}
	// fetch environment
	if err := p.PlatformSH.Fetch(); err != nil {
		return tracerr.Wrap(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return tracerr.Errorf("environment '%s' not found", envName)
	}
	// fetch variables
	done := output.Duration("Fetch variables.")
	vars, err := p.PlatformSH.Variables(env)
	if err != nil {
		return tracerr.Wrap(err)
	}
	for k, v := range vars {
		if err := p.VarSet(k, v); err != nil {
			return tracerr.Wrap(err)
		}
	}
	pvars, err := p.PlatformSH.PlatformVariables(env, p.Apps[0].Name)
	if err != nil {
		return tracerr.Wrap(err)
	}
	for k, v := range pvars {
		if err := p.VarSet(k, def.InterfaceToString(v)); err != nil {
			return tracerr.Wrap(err)
		}
	}

	if err := p.Save(); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// PlatformSHSync syncs the given platform.sh environment to the current local project.
func (p *Project) PlatformSHSync(envName string) error {

	done := output.Duration(fmt.Sprintf("Sync %s environment.", envName))

	p.PlatformSHSyncVariables(envName)

	done()
	return nil
}
