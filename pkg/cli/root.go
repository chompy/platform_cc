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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/pkg/output"
)

// RootCmd is the top level command.
var RootCmd = &cobra.Command{
	Use:     "platform_cc [-v verbose]",
	Version: "",
	Run: func(cmd *cobra.Command, args []string) {
		commandIntro(cmd.Version)
		output.WriteStdout("\nAvailable Commands:\n")
		displayCommandList(cmd)
	},
}

// Execute - run root command
func Execute() error {
	// hack that allows old style semicolon (:) seperated
	// subcommands to work
	args := make([]string, 1)
	args[0] = os.Args[0]
	if len(os.Args) > 1 {
		args = append(args, strings.Split(os.Args[1], ":")...)
		args = append(args, os.Args[2:]...)
	}
	os.Args = args
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "show more verbose output")
}
