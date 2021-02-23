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

package output

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CommandIntro displays introduction information about Platform.CC
func CommandIntro(version string) {
	fmt.Println(colorSuccess(strings.Repeat("=", 32)))
	fmt.Println(" PLATFORM.CC BY CONTEXTUAL CODE ")
	fmt.Println(colorSuccess(strings.Repeat("=", 32)))
	fmt.Println(color(" VERSION "+version, 35))
}

// CommandList displays recursive list of available commands.
func CommandList(cmd *cobra.Command) {
	var ittListCmd func(cmd *cobra.Command, prefix string)
	ittListCmd = func(cmd *cobra.Command, prefix string) {
		if cmd.HasSubCommands() {
			for _, scmd := range cmd.Commands() {
				ittListCmd(scmd, prefix+":"+scmd.Name())
			}
			return
		}
		name := "  " + prefix
		fmt.Println(colorSuccess(name) + strings.Repeat(" ", 42-len(name)) + cmd.Short)
	}
	for _, scmd := range cmd.Commands() {
		if !scmd.HasSubCommands() {
			fmt.Println(colorSuccess("  " + scmd.Name()))
		}
	}
	for _, scmd := range cmd.Commands() {
		if scmd.HasSubCommands() {
			fmt.Print(color(scmd.Name(), 94))
			if len(scmd.Aliases) > 0 {
				fmt.Print(color(" ["+strings.Join(scmd.Aliases, ",")+"]", 35))
			}
			fmt.Print("\n")
			ittListCmd(scmd, scmd.Name())
		}
	}
}
