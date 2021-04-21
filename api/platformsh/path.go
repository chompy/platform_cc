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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vaughan0/go-ini"
	"github.com/ztrue/tracerr"
)

const platformShGitPattern = `^([a-z0-9]{12,})@git\.(([a-z0-9\-]+\.)?platform\.sh):([a-z0-9]{12,})\.git$`
const gitConfigPath = ".git/config"

// parseProjectGit returns Platform.sh project id and host name from Git remote.
func parseProjectGit(path string) (string, string, error) {
	projGitConfigPath := filepath.Join(path, gitConfigPath)
	// load git config file
	conf, err := ini.LoadFile(projGitConfigPath)
	if err != nil {
		return "", "", tracerr.Wrap(err)
	}
	// compile regexp
	gitMatchRegex, err := regexp.Compile(platformShGitPattern)
	if err != nil {
		return "", "", tracerr.Wrap(err)
	}
	// find git remote
	for name, section := range conf {
		if strings.HasPrefix(name, "remote") {
			res := gitMatchRegex.FindAllStringSubmatch(section["url"], -1)
			if len(res) > 0 && len(res[0]) > 2 {
				return res[0][1], res[0][2], nil
			}
		}
	}
	return "", "", fmt.Errorf("platform.sh git remote url not found")
}

// FindRoot returns the root directory for the Platform.sh project.
func FindRoot(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	pathSplit := strings.Split(path, string(os.PathSeparator))
	if pathSplit[0] == "" {
		pathSplit[0] = string(os.PathSeparator)
	}
	for i := range pathSplit {
		currentPath := filepath.Join(pathSplit[0 : i+1]...)
		_, _, err := parseProjectGit(currentPath)
		if err == nil {
			return currentPath, nil
		}
	}
	return "", fmt.Errorf("could not find valid platform.sh project in current path")
}
