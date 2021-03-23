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

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/spf13/cobra"
)

const termColorCommand = 32
const termColorTopLevelCommand = 94
const termColorCommandAlias = 35

type flatCommand struct {
	Command *cobra.Command
	Level   int
}

func (f flatCommand) Prefix() []string {
	cmdItt := f.Command
	prefix := make([]string, 0)
	for {
		cmdItt = cmdItt.Parent()
		if !cmdItt.HasParent() {
			break
		}
		prefix = append(prefix, cmdItt.Name())
	}
	return prefix
}

func (f flatCommand) ParentCommand() *cobra.Command {
	cmdItt := f.Command
	outCmd := cmdItt
	for {
		cmdItt = cmdItt.Parent()
		if !cmdItt.HasParent() {
			break
		}
		outCmd = cmdItt
	}
	return outCmd
}

func (f flatCommand) IsMatch(name string) bool {
	if name == f.Command.Name() {
		return true
	}
	for _, alias := range f.Command.Aliases {
		if alias == name {
			return true
		}
	}
	return false
}

func (f flatCommand) PrefixMatch(prefix string) bool {
	if strings.HasPrefix(f.Command.Name(), prefix) {
		return true
	}
	for _, alias := range f.Command.Aliases {
		if strings.HasPrefix(alias, prefix) {
			return true
		}
	}
	return false
}

// flatCommandList returns the entire command tree flattened in to a single list.
func flatCommandList(cmd *cobra.Command) []flatCommand {
	out := make([]flatCommand, 0)
	var ittListCmd func(cmd *cobra.Command, level int)
	ittListCmd = func(cmd *cobra.Command, level int) {
		if cmd.Hidden {
			return
		}
		if cmd.HasSubCommands() {
			for _, scmd := range cmd.Commands() {
				ittListCmd(scmd, level+1)
			}
			return
		}
		if level > 1 {
			out = append(
				out,
				flatCommand{
					Command: cmd,
					Level:   level,
				},
			)
		}
	}
	for _, scmd := range cmd.Commands() {
		if scmd.HasSubCommands() || scmd.Hidden {
			continue
		}
		out = append(
			out,
			flatCommand{
				Command: scmd,
				Level:   2,
			},
		)
	}
	ittListCmd(cmd, 0)
	return out
}

func flatCommandListFilter(cmd *cobra.Command, name string) []flatCommand {
	if !cmd.HasSubCommands() {
		return []flatCommand{}
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return flatCommandList(cmd)
	}
	args := strings.Split(name, ":")
	for _, scmd := range cmd.Commands() {
		fc := flatCommand{Command: scmd}
		if !scmd.HasSubCommands() || !fc.IsMatch(args[0]) {
			continue
		}
		return flatCommandListFilter(scmd, strings.Join(args[1:], ":"))
	}
	if len(args) > 1 {
		return []flatCommand{}
	}
	out := make([]flatCommand, 0)
	for _, scmd := range cmd.Commands() {
		fc := flatCommand{Command: scmd}
		if !fc.PrefixMatch(args[0]) {
			continue
		}
		out = append(out, fc)
	}
	return out
}

// displayCommandList prints a list of commands to the terminal.
func displayCommandList(cmd *cobra.Command) {
	cmdList := flatCommandList(cmd)
	topLevelCmdName := ""
	for _, flatCmd := range cmdList {
		prevTopLevelCmdName := topLevelCmdName
		tabLength := 42
		topCmd := flatCmd.ParentCommand()
		topLevelCmdName = topCmd.Name()
		topLevelAliases := topCmd.Aliases
		// display top level command with aliases
		if prevTopLevelCmdName != topLevelCmdName && flatCmd.Command != topCmd {
			output.WriteStdout(
				output.Color(topLevelCmdName, termColorTopLevelCommand),
			)
			if len(topLevelAliases) > 0 {
				output.WriteStdout(output.Color(
					" ["+strings.Join(topLevelAliases, ",")+"]", termColorCommandAlias,
				))
			}
			output.WriteStdout("\n")
		}
		// tab over sub commands
		if flatCmd.Level > 1 {
			tabLength -= 2
			output.WriteStdout("  ")
		}
		// display sub command name
		printName := strings.Join(flatCmd.Prefix(), ":")
		if printName != "" {
			printName += ":"
		}
		printName += flatCmd.Command.Name()
		tabLength -= len(printName)
		output.WriteStdout(
			output.Color(printName, termColorCommand) +
				strings.Repeat(" ", tabLength) + flatCmd.Command.Short + "\n",
		)
	}
}

// ListCmd lists all available commands.
var ListCmd = &cobra.Command{
	Use:     "list",
	Version: "",
	Run: func(cmd *cobra.Command, args []string) {
		commandIntro(RootCmd.Version)
		output.WriteStdout("\nAvailable Commands:\n")
		displayCommandList(RootCmd)
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
