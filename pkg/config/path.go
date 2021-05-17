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
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

const configPerm = 0766
const userConfigPath = "~/.config/platformcc"

func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	usr, err := user.Current()
	if err != nil {
		return path
	}
	return filepath.Join(usr.HomeDir, path[1:])
}

func initConfig() error {
	err := os.MkdirAll(Path(), configPerm)
	if os.IsExist(err) {
		return nil
	}
	return errors.WithStack(err)
}

func pathTo(name string) string {
	return filepath.Join(Path(), name)
}

// Path returns the path to the config directory.
func Path() string {
	return expandPath(userConfigPath)
}
