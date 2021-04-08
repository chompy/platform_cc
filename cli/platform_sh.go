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
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gitlab.com/contextualcode/platform_cc/api/platformsh"

	"github.com/spf13/cobra"
	"github.com/ztrue/tracerr"
)

var platformShCmd = &cobra.Command{
	Use:     "platform-sh",
	Aliases: []string{"platform_sh", "platformsh", "psh"},
	Short:   "Manage Platform.sh options.",
}

var platformShLoginCmd = &cobra.Command{
	Use: "login",
	Run: func(cmd *cobra.Command, args []string) {
		handleError(platformsh.APILogin())
	},
}

var platformShSshCmd = &cobra.Command{
	Use: "ssh environment [-s service]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			handleError(tracerr.Errorf("environment not provided"))
		}
		// fetch project
		proj, err := getProject(true)
		handleError(err)
		if proj.PlatformSH == nil || proj.PlatformSH.ID == "" {
			handleError(tracerr.Errorf("platform.sh project not found"))
		}
		// get def
		def, err := getDefFromCommand(cmd, proj)
		handleError(err)
		// fetch environments
		handleError(proj.PlatformSH.FetchEnvironments())
		// get psh environment
		env := proj.PlatformSH.GetEnvironment(args[0])
		if env == nil {
			handleError(tracerr.Errorf("cannot find environment %s", args[0]))
		}
		// fetch ssh url
		output.WriteStdout(
			proj.PlatformSH.SSHUrl(env, proj.GetDefinitionName(def)),
		)
	},
}

var platformShSyncCmd = &cobra.Command{
	Use: "sync environment [--skip-variables] [--skip-mounts] [--skip-databases]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			handleError(tracerr.Errorf("environment not provided"))
		}
		proj, err := getProject(true)
		handleError(err)
		if !checkFlag(cmd, "skip-variables") {
			handleError(proj.PlatformSHSyncVariables(args[0]))
		}
		if !checkFlag(cmd, "skip-mounts") || !checkFlag(cmd, "skip-databases") {
			handleError(proj.Start())
		}
		if !checkFlag(cmd, "skip-mounts") {
			handleError(proj.PlatformSHSyncMounts(args[0]))
		}
		if !checkFlag(cmd, "skip-databases") {
			handleError(proj.PlatformSHSyncDatabases(args[0]))
		}
		if !checkFlag(cmd, "skip-mounts") || !checkFlag(cmd, "skip-databases") {
			handleError(proj.Stop())
		}
	},
}

func init() {
	platformShSshCmd.PersistentFlags().StringP("service", "s", "", "name of service/application/worker")
	platformShCmd.AddCommand(platformShLoginCmd)
	platformShCmd.AddCommand(platformShSshCmd)
	platformShSyncCmd.Flags().Bool("skip-variables", false, "Skip variable sync.")
	platformShSyncCmd.Flags().Bool("skip-mounts", false, "Skip mount sync.")
	platformShSyncCmd.Flags().Bool("skip-databases", false, "Skip database sync.")
	platformShCmd.AddCommand(platformShSyncCmd)
	RootCmd.AddCommand(platformShCmd)
}
