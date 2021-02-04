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
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

var appCmd = &cobra.Command{
	Use:     "application [-n name]",
	Aliases: []string{"app"},
	Short:   "Manage applications.",
}

var appBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run application build commands.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		app, err := getApp(appCmd, proj)
		handleError(err)
		c := proj.NewContainer(app)
		handleError(c.Build())
		handleError(c.Commit())
	},
}

var appDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run application deploy hook.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		app, err := getApp(appCmd, proj)
		handleError(err)
		c := proj.NewContainer(app)
		handleError(c.Deploy())
	},
}

var appShellCmd = &cobra.Command{
	Use:     "shell [--root] [command]",
	Aliases: []string{"sh"},
	Short:   "Shell in to an application.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		app, err := getApp(appCmd, proj)
		handleError(err)
		user := "web"
		if cmd.PersistentFlags().Lookup("root").Value.String() == "true" {
			user = "root"
		}
		c := proj.NewContainer(app)
		shellCmd := []string{}
		if len(args) > 0 {
			shellCmd = []string{
				"sh", "-c", strings.Join(args, " "),
			}
		}
		handleError(c.Shell(user, shellCmd))
	},
}

var appCommitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"cmt", "cm", "c"},
	Short:   "Commit container state.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		app, err := getApp(appCmd, proj)
		handleError(err)
		c := proj.NewContainer(app)
		handleError(c.Commit())
	},
}

var appDeleteCommitCmd = &cobra.Command{
	Use:     "delete_commit",
	Aliases: []string{"dc", "dcmt", "dcommit", "dcm"},
	Short:   "Delete container comit.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		app, err := getApp(appCmd, proj)
		handleError(err)
		c := proj.NewContainer(app)
		handleError(c.DeleteCommit())
	},
}

var appLogsCmd = &cobra.Command{
	Use:   "logs [-f follow]",
	Short: "Display logs for application.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		app, err := getApp(appCmd, proj)
		handleError(err)
		handleError(proj.NewContainer(app).LogStdout())
		time.Sleep(time.Second)
		// follow logs
		followFlag := cmd.Flags().Lookup("follow")
		if followFlag != nil && followFlag.Value.String() != "false" {
			select {}
		}
	},
}

func init() {
	appShellCmd.PersistentFlags().Bool("root", false, "shell as root")
	appCmd.PersistentFlags().StringP("name", "n", "", "name of application")
	appLogsCmd.Flags().BoolP("follow", "f", false, "follow logs")
	appCmd.AddCommand(appBuildCmd)
	appCmd.AddCommand(appDeployCmd)
	appCmd.AddCommand(appShellCmd)
	appCmd.AddCommand(appCommitCmd)
	appCmd.AddCommand(appDeleteCommitCmd)
	appCmd.AddCommand(appLogsCmd)
	RootCmd.AddCommand(appCmd)
}
