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

	"github.com/spf13/cobra"
	"github.com/ztrue/tracerr"
)

var platformShCmd = &cobra.Command{
	Use:     "platform-sh",
	Aliases: []string{"platform_sh", "platformsh", "psh"},
	Short:   "Manage Platform.sh options.",
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
		def, err := getDefFromCommand(containerCmd, proj)
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
	Use: "sync environment [--skip-variables] [--skip-mounts] [--skip-services]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			handleError(tracerr.Errorf("environment not provided"))
		}
		proj, err := getProject(true)
		handleError(err)
		if !checkFlag(cmd, "skip-variables") {
			handleError(proj.PlatformSHSyncVariables(args[0]))
		}
	},
}

func init() {
	platformShCmd.AddCommand(platformShSshCmd)
	platformShSyncCmd.Flags().Bool("skip-variables", false, "Skip variable sync.")
	platformShCmd.AddCommand(platformShSyncCmd)
	RootCmd.AddCommand(platformShCmd)
}
