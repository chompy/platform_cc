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
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
)

// Project contains information about a Platform.sh project.
type Project struct {
	LocalPath    string
	Host         string
	ID           string
	Title        string
	Environments []Environment
	apiURL       string
}

// LoadProjectFromPath returns information about a Platform.sh from a path.
func LoadProjectFromPath(path string) (*Project, error) {
	path, err := FindRoot(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pid, host, err := parseProjectGit(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	output.Info(fmt.Sprintf("Found Platform.sh project '%s.'", pid))
	return &Project{
		LocalPath: path,
		Host:      host,
		ID:        pid,
		apiURL:    apiURL,
	}, nil
}
