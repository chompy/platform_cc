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
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/pkg/output"
)

// VarSet sets a project variable.
func (p *Project) VarSet(key string, value string) error {
	output.Info(
		fmt.Sprintf("Set var '%s.'", key),
	)
	value = strings.TrimSpace(value)
	output.LogDebug(fmt.Sprintf("Set var '%s.'", key), value)
	return errors.WithStack(p.Variables.Set(key, value))
}

// VarGet retrieves a project variable.
func (p *Project) VarGet(key string) (string, error) {
	output.Info(
		fmt.Sprintf("Get var '%s.'", key),
	)
	out := p.Variables.GetString(key)
	output.LogDebug(fmt.Sprintf("Get var '%s.'", key), out)
	return out, nil
}

// VarDelete deletes a project variable.
func (p *Project) VarDelete(key string) error {
	output.Info(
		fmt.Sprintf("Delete var '%s.'", key),
	)
	p.Variables.Delete(key)
	return nil
}
