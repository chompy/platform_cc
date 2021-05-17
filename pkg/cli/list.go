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
	"sort"
	"strings"

	"gitlab.com/contextualcode/platform_cc/pkg/output"

	"github.com/spf13/cobra"
)

const listDescSpacing = 42
const termColorCommand = 32
const termColorTopLevelCommand = 94
const termColorCommandAlias = 35

// listCommands returns a list of every command.
func listCommands(cmd *cobra.Command) []string {
	out := make([]string, 0)
	for _, child := range cmd.Commands() {
		if child.Hidden {
			continue
		}
		if !child.HasSubCommands() {
			out = append(out, child.Name())
			continue
		}
		for _, childChildCmd := range listCommands(child) {
			out = append(out, child.Name()+":"+childChildCmd)
		}
	}
	sort.Strings(out)
	return out
}

// displayCommandList prints a list of commands to the terminal.
func displayCommandList(cmd *cobra.Command) {
	cmdList := listCommands(cmd)
	topLevelCmdName := ""
	// helper function to write command details to term
	printCommand := func(name string) {
		currentCmd, _, err := cmd.Find(strings.Split(name, ":"))
		if err != nil {
			return
		}
		tabLength := listDescSpacing
		tabLength -= len(name)
		output.WriteStdout(
			strings.Repeat(" ", 2) +
				output.Color(name, termColorCommand) +
				strings.Repeat(" ", tabLength) + currentCmd.Short + "\n",
		)
	}
	// top level commands
	for _, name := range cmdList {
		if !strings.Contains(name, ":") {
			printCommand(name)
		}
	}
	// child commands
	for _, name := range cmdList {
		if !strings.Contains(name, ":") {
			continue
		}
		// display top level command to categorize sub commands
		currentTopLevelCmd := strings.Split(name, ":")[0]
		if currentTopLevelCmd != topLevelCmdName {
			topLevelCmdName = currentTopLevelCmd
			topLevelCmd, _, err := cmd.Find(strings.Split(topLevelCmdName, ":"))
			if err != nil {
				continue
			}
			output.WriteStdout(
				output.Color(topLevelCmdName, termColorTopLevelCommand),
			)
			if len(topLevelCmd.Aliases) > 0 {
				output.WriteStdout(output.Color(
					" ["+strings.Join(topLevelCmd.Aliases, ",")+"]", termColorCommandAlias,
				))
			}
			output.WriteStdout("\n")
		}
		printCommand(name)
	}
}

// ListCmd lists all available commands.
var ListCmd = &cobra.Command{
	Use:     "list",
	Version: "",
	Short:   "Show this list.",
	Run: func(cmd *cobra.Command, args []string) {
		commandIntro(RootCmd.Version)
		output.WriteStdout("\nAvailable Commands:\n")
		displayCommandList(RootCmd)
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
