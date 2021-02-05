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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gitlab.com/contextualcode/platform_cc/api/router"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj", "p"},
	Short:   "Manage project in current working directory.",
}

var projectStartCmd = &cobra.Command{
	Use:   "start [--no-build] [--no-router] [--no-commit] [-s slot]",
	Short: "Start a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		// determine volume slot
		slot, err := strconv.Atoi(cmd.Flags().Lookup("slot").Value.String())
		handleError(err)
		proj.SetVolumeSlot(slot)
		// start
		handleError(proj.Start())
		noCommitFlag := cmd.Flags().Lookup("no-commit")
		noBuildFlag := cmd.Flags().Lookup("no-build")
		if noBuildFlag == nil || noBuildFlag.Value.String() == "false" {
			doCommit := noCommitFlag == nil || noCommitFlag.Value.String() == "false"
			handleError(proj.Build(false, doCommit))
		}
		noRouterFlag := cmd.Flags().Lookup("no-router")
		if noRouterFlag == nil || noRouterFlag.Value.String() == "false" {
			handleError(router.Start())
			handleError(router.AddProjectRoutes(proj))
		}

	},
}

var projectStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		handleError(proj.Stop())
		handleError(router.DeleteProjectRoutes(proj))
	},
}

var projectRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.Stop())
		handleError(proj.Start())
	},
}

var projectPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull all Docker images for project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.Pull())
	},
}

var projectBuildCmd = &cobra.Command{
	Use:   "build [--no-commit]",
	Short: "Build a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		noCommitFlag := cmd.Flags().Lookup("no-commit")
		doCommit := noCommitFlag == nil || noCommitFlag.Value.String() != "false"
		handleError(proj.Build(true, doCommit))
	},
}

var projectDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run deploy hooks for project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.Deploy())
	},
}

var projectPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		handleError(proj.Purge())
	},
}

var projectConfigJSONCmd = &cobra.Command{
	Use:   "configjson",
	Short: "Dump config.json.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		configJSON, err := proj.BuildConfigJSON(proj.Apps[0])
		handleError(err)
		fmt.Println(string(configJSON))
	},
}

var projectStatusCmd = &cobra.Command{
	Use:   "status [--json]",
	Short: "Display status of project containers.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		status := proj.Status()
		// json out
		jsonFlag := cmd.Flags().Lookup("json")
		if jsonFlag != nil && jsonFlag.Value.String() != "false" {
			out, err := json.Marshal(status)
			handleError(err)
			fmt.Println(string(out))
			return
		}
		// table out
		data := make([][]string, 0)
		for _, s := range status {
			runningStr := "stopped"
			ipAddrStr := "n/a"
			if s.Running {
				runningStr = "running"
			}
			if s.IPAddress != "" {
				ipAddrStr = s.IPAddress
			}
			data = append(data, []string{
				s.Name,
				s.Type,
				runningStr,
				ipAddrStr,
			})
		}
		drawTable(
			[]string{"Name", "Type", "Status", "IP Address"},
			data,
		)
	},
}

var projectLogsCmd = &cobra.Command{
	Use:   "logs [-f follow]",
	Short: "Display logs for all project containers.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		for _, app := range proj.Apps {
			handleError(proj.NewContainer(app).LogStdout())
		}
		for _, service := range proj.Services {
			handleError(proj.NewContainer(service).LogStdout())
		}
		time.Sleep(time.Second)
		// follow logs
		followFlag := cmd.Flags().Lookup("follow")
		if followFlag != nil && followFlag.Value.String() != "false" {
			select {}
		}
	},
}

func init() {
	projectStartCmd.Flags().Bool("no-build", false, "skip building project")
	projectStartCmd.Flags().Bool("no-router", false, "skip adding routes to router")
	projectStartCmd.Flags().Bool("no-commit", false, "don't commit the container after being built")
	projectStartCmd.Flags().IntP("slot", "s", 0, "set volume slot")
	projectBuildCmd.Flags().Bool("no-commit", false, "don't commit the container after being built")
	projectStatusCmd.Flags().Bool("json", false, "JSON output")
	projectLogsCmd.Flags().BoolP("follow", "f", false, "follow logs")
	projectCmd.AddCommand(projectStartCmd)
	projectCmd.AddCommand(projectStopCmd)
	projectCmd.AddCommand(projectRestartCmd)
	projectCmd.AddCommand(projectPullCmd)
	projectCmd.AddCommand(projectBuildCmd)
	projectCmd.AddCommand(projectDeployCmd)
	projectCmd.AddCommand(projectPurgeCmd)
	projectCmd.AddCommand(projectConfigJSONCmd)
	projectCmd.AddCommand(projectStatusCmd)
	projectCmd.AddCommand(projectLogsCmd)
	RootCmd.AddCommand(projectCmd)
}
