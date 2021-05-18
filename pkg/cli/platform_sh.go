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
	"fmt"

	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/platformsh"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/project"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var platformShCmd = &cobra.Command{
	Use:     "platform-sh",
	Aliases: []string{"platform_sh", "platformsh", "psh"},
	Short:   "Manage Platform.sh options.",
}

var platformShLoginCmd = &cobra.Command{
	Use: "login",
	Run: func(cmd *cobra.Command, args []string) {
		handleError(platformsh.Login())
	},
}

var platformShSSHCmd = &cobra.Command{
	Use: "ssh [-e environment] [-s service] [--pipe]",
	Run: func(cmd *cobra.Command, args []string) {
		// fetch project
		proj, err := getProject(true)
		handleError(err)
		if proj.PlatformSH == nil || proj.PlatformSH.ID == "" {
			handleError(errors.WithStack(platformsh.ErrProjectNotFound))
		}
		// get def
		def, err := getDefFromCommand(cmd, proj)
		handleError(err)
		// fetch environments
		handleError(proj.PlatformSH.FetchEnvironments())
		// get psh environment
		env, err := getPlatformShEnvironment(cmd, proj)
		handleError(err)
		// dump ssh url to term if pipe option provided
		if checkFlag(cmd, "pipe") {
			output.WriteStdout(
				proj.PlatformSH.SSHUrl(env, proj.GetDefinitionName(def)) + "\n",
			)
			return
		}
		// tty
		handleError(proj.PlatformSH.SSHTerminal(env, proj.GetDefinitionName(def)))
	},
}

var platformShSyncCmd = &cobra.Command{
	Use: "sync [-e environment] [--skip-variables] [--skip-mounts] [--skip-databases]",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		// get psh environment
		env, err := getPlatformShEnvironment(cmd, proj)
		handleError(err)
		output.Info(fmt.Sprintf("Sync from %s (%s) environment.", env.Name, env.MachineName))
		// set mount strat to volume
		output.Info("Set mount strategy to volume.")
		proj.Options[project.OptionMountStrategy] = project.MountStrategyVolume
		handleError(proj.Save())
		// perform sync tasks
		if !checkFlag(cmd, "skip-variables") {
			handleError(proj.PlatformSHSyncVariables(env.Name))
		}
		if !checkFlag(cmd, "skip-mounts") || !checkFlag(cmd, "skip-databases") {
			handleError(proj.Start())
		}
		if !checkFlag(cmd, "skip-mounts") {
			handleError(proj.PlatformSHSyncMounts(env.Name))
		}
		if !checkFlag(cmd, "skip-databases") {
			handleError(proj.PlatformSHSyncDatabases(env.Name))
		}
		if !checkFlag(cmd, "skip-mounts") || !checkFlag(cmd, "skip-databases") {
			handleError(proj.Stop())
		}
	},
}

func init() {
	platformShSSHCmd.PersistentFlags().StringP("service", "s", "", "name of service/application/worker")
	platformShSSHCmd.PersistentFlags().StringP("environment", "e", "", "name of environment")
	platformShSSHCmd.Flags().Bool("pipe", false, "return ssh url instead of creating interactive terminal")
	platformShCmd.AddCommand(platformShLoginCmd)
	platformShCmd.AddCommand(platformShSSHCmd)
	platformShSyncCmd.PersistentFlags().StringP("environment", "e", "", "name of environment")
	platformShSyncCmd.Flags().Bool("skip-variables", false, "Skip variable sync.")
	platformShSyncCmd.Flags().Bool("skip-mounts", false, "Skip mount sync.")
	platformShSyncCmd.Flags().Bool("skip-databases", false, "Skip database sync.")
	platformShCmd.AddCommand(platformShSyncCmd)
	RootCmd.AddCommand(platformShCmd)
}
