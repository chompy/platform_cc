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
	"log"
	"strings"

	"github.com/ztrue/tracerr"
)

// VarSet sets a project variable.
func (p *Project) VarSet(key string, value string) error {
	log.Printf("Set var '%s.'", key)
	keySplit := strings.Split(key, ":")
	if len(keySplit) != 2 {
		return tracerr.Wrap(fmt.Errorf("invalid variable key"))
	}
	if p.Variables[keySplit[0]] == nil {
		p.Variables[keySplit[0]] = make(map[string]string)
	}
	p.Variables[keySplit[0]][keySplit[1]] = value
	return nil
}

// VarGet retrieves a project variable.
func (p *Project) VarGet(key string) (string, error) {
	keySplit := strings.Split(key, ":")
	if len(keySplit) != 2 {
		return "", tracerr.Wrap(fmt.Errorf("invalid variable key"))
	}
	if p.Variables[keySplit[0]] == nil {
		return "", nil
	}
	return p.Variables[keySplit[0]][keySplit[1]], nil
}

// VarDelete deletes a project variable.
func (p *Project) VarDelete(key string) error {
	log.Printf("Delete var '%s.'", key)
	keySplit := strings.Split(key, ":")
	if len(keySplit) != 2 {
		return tracerr.Wrap(fmt.Errorf("invalid variable key"))
	}
	if p.Variables[keySplit[0]] == nil {
		return nil
	}
	delete(p.Variables[keySplit[0]], keySplit[1])
	return nil
}
