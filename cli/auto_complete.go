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

package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// listCommandAliases returns a list of every possible command combination with aliases.
func listCommandAliases(cmd *cobra.Command) []string {
	out := make([]string, 0)
	for _, child := range cmd.Commands() {
		aliases := child.Aliases
		aliases = append(aliases, child.Name())
		if !child.HasSubCommands() {
			out = append(out, aliases...)
			continue
		}
		for _, alias := range aliases {
			childChildCmds := listCommandAliases(child)
			for _, childChildCmd := range childChildCmds {
				out = append(out, alias+":"+childChildCmd)
			}
		}
	}
	return out
}

func filterListCommandAliases(cmd *cobra.Command, filter string) []string {
	allCommands := listCommandAliases(cmd)
	filterArgSplit := strings.Split(filter, ":")
	out := make([]string, 0)
	for _, cmdStr := range allCommands {
		if cmdStr == filter {
			return []string{}
		} else if filter == "" || strings.HasPrefix(cmdStr, filter) {
			// if a completed parent command has already been typed out then don't
			// include it in results
			if len(filterArgSplit) > 1 {
				for _, arg := range filterArgSplit[0 : len(filterArgSplit)-1] {
					cmdStr = strings.TrimPrefix(cmdStr, arg+":")
				}
			}
			// check for duplicate before adding
			hasOut := false
			for _, v := range out {
				if v == cmdStr {
					hasOut = true
					break
				}
			}
			if !hasOut {
				out = append(out, cmdStr)
			}
		}
	}
	return out
}

// AutoCompleteListCmd list every possible command for Bash auto-complete.
var AutoCompleteListCmd = &cobra.Command{
	Hidden: true,
	Use:    "_ac",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		output.Logging = false
		// format args so that the command is only one arg
		// bash seems to treat ":" as a seperate arg
		if len(args) > 0 {
			argStr := ""
			for _, arg := range strings.Split(args[0], ":") {
				argStr += strings.TrimSpace(arg) + ":"
			}
			argStr = strings.TrimSuffix(argStr, ":")
			args = strings.Split(argStr, " ")
		}
		// only perform command auto complete if arg count is one or less
		if len(args) <= 1 {
			out := filterListCommandAliases(RootCmd, args[0])
			output.WriteStdout(strings.Join(out, " "))
		}
	},
}

func init() {
	RootCmd.AddCommand(AutoCompleteListCmd)
}
