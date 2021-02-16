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

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

var containerCmd = &cobra.Command{
	Use:     "container [-n name]",
	Aliases: []string{"cont", "cnt", "application", "ap", "app", "appl"},
	Short:   "Manage individual containers for applications, services, and workers.",
}

var containerAppDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run application deploy hook.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDef(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.Deploy())
				return
			}
		}
		handleError(fmt.Errorf("can only run deploy hooks on applications"))
	},
}

var containerShellCmd = &cobra.Command{
	Use:     "shell [--root] [command]",
	Aliases: []string{"sh"},
	Short:   "Shell in to a container.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDef(containerCmd, proj)
		handleError(err)
		user := "root"
		switch d.(type) {
		case def.App:
			{
				user = "web"
				if cmd.PersistentFlags().Lookup("root").Value.String() == "true" {
					user = "root"
				}
				break
			}
		}
		c := proj.NewContainer(d)
		shellCmd := []string{}
		if len(args) > 0 {
			shellCmd = []string{
				"bash", "--login", "-c", strings.Join(args, " "),
			}
		}
		handleError(c.Shell(user, shellCmd))
	},
}

var containerAppCommitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"cmt", "cm", "c"},
	Short:   "Commit container state.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDef(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.Commit())
				return
			}
		}
		handleError(fmt.Errorf("can only commit applications"))
	},
}

var containerAppDeleteCommitCmd = &cobra.Command{
	Use:     "delete_commit",
	Aliases: []string{"dc", "dcmt", "dcommit", "dcm"},
	Short:   "Delete container comit.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDef(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.DeleteCommit())
				return
			}
		}
		handleError(fmt.Errorf("can only commit applications"))
	},
}

var containerLogsCmd = &cobra.Command{
	Use:   "logs [-f follow]",
	Short: "Display logs for container.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		followFlag := cmd.Flags().Lookup("follow")
		hasFollow := followFlag != nil && followFlag.Value.String() != "false"
		d, err := getDef(containerCmd, proj)
		handleError(err)
		handleError(proj.NewContainer(d).LogStdout(hasFollow))
		if hasFollow {
			select {}
		}
	},
}

func init() {
	containerShellCmd.PersistentFlags().Bool("root", false, "shell as root")
	containerCmd.PersistentFlags().StringP("name", "n", "", "name of application")
	containerLogsCmd.Flags().BoolP("follow", "f", false, "follow logs")
	containerCmd.AddCommand(containerAppDeployCmd)
	containerCmd.AddCommand(containerShellCmd)
	containerCmd.AddCommand(containerAppCommitCmd)
	containerCmd.AddCommand(containerAppDeleteCommitCmd)
	containerCmd.AddCommand(containerLogsCmd)
	RootCmd.AddCommand(containerCmd)
}
