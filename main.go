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

package main

import (
	"gitlab.com/contextualcode/platform_cc/v2/pkg/cli"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
)

var version string

func main() {
	output.Enable = true
	output.Logging = true
	if err := output.LogRotate(); err != nil {
		output.Warn("Error while trimming log, " + err.Error())
	}
	if version == "" {
		version = "DEV"
	}
	cli.RootCmd.Version = version
	cli.Execute()
}
