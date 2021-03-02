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
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gitlab.com/contextualcode/platform_cc/api/router"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj", "p"},
	Short:   "Manage project in current working directory.",
}

var projectStartCmd = &cobra.Command{
	Use:   "start [--rebuild] [--no-build] [--no-router] [--no-commit] [--no-validate] [-s slot]",
	Short: "Start a project.",
	Run: func(cmd *cobra.Command, args []string) {
		projectStart(cmd, nil, -1)
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
	Use:   "restart [--rebuild] [--no-build] [--no-commit] [--no-validate]",
	Short: "Restart a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		// get status and retrieve slot + stop
		slot := -1
		status := proj.Status()
		if len(status) > 0 {
			slot = status[0].Slot
			handleError(proj.Stop())
		}
		// remove from router
		handleError(router.DeleteProjectRoutes(proj))
		// start
		projectStart(cmd, proj, slot)
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

var projectDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run deploy hooks for project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.Deploy())
	},
}

var projectPostDeployCmd = &cobra.Command{
	Use:     "post-deploy",
	Aliases: []string{"postdeploy", "pd"},
	Short:   "Run post-deploy hooks for project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.PostDeploy())
	},
}

var projectPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		handleError(proj.Purge())
		handleError(router.DeleteProjectRoutes(proj))
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
		output.WriteStdout(string(configJSON) + "\n")
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
			output.WriteStdout(string(out) + "\n")
			return
		}
		// table out
		data := make([][]string, 0)
		slot := "n/a"
		activeCount := 0
		for _, s := range status {
			ipAddrStr := "n/a"
			if s.Running {
				activeCount++
			}
			if s.IPAddress != "" {
				ipAddrStr = s.IPAddress
			}
			if s.Slot > 0 && slot == "n/a" {
				slot = fmt.Sprintf("%d", s.Slot)
			}
			serviceType := s.Type
			if s.Committed {
				serviceType = "[c] " + s.Type
			}

			data = append(data, []string{
				fmt.Sprintf("[%s] %s", string(s.ObjectType), s.Name),
				serviceType,
				s.State,
				ipAddrStr,
			})
		}
		log.Println(data)
		output.WriteStdout("\n")
		output.WriteStdout(fmt.Sprintf("ID\t\t%s\n", proj.ID))
		output.WriteStdout(fmt.Sprintf("SLOT\t\t%s\n", slot))
		output.WriteStdout(fmt.Sprintf("ACTIVE\t\t%d/%d\n", activeCount, len(status)))
		drawTable(
			[]string{"Name", "Type", "Status", "IP Address"},
			data,
		)
		drawKeys()
		output.WriteStdout("\n")
	},
}

var projectLogsCmd = &cobra.Command{
	Use:   "logs [-f follow]",
	Short: "Display logs for all project containers.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		followFlag := cmd.Flags().Lookup("follow")
		hasFollow := followFlag != nil && followFlag.Value.String() != "false"
		for _, app := range proj.Apps {
			if hasFollow {
				go proj.NewContainer(app).LogStdout(hasFollow)
			} else {
				handleError(proj.NewContainer(app).LogStdout(hasFollow))
			}
		}
		for _, service := range proj.Services {
			if hasFollow {
				go proj.NewContainer(service).LogStdout(hasFollow)
			} else {
				handleError(proj.NewContainer(service).LogStdout(hasFollow))
			}
		}
		if hasFollow {
			select {}
		}
	},
}

var projectRoutesCmd = &cobra.Command{
	Use:     "routes [--json]",
	Aliases: []string{"route", "rout", "ro"},
	Short:   "List routes for project.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		// json out
		jsonFlag := cmd.Flags().Lookup("json")
		if jsonFlag != nil && jsonFlag.Value.String() != "false" {
			routes := proj.Routes
			for i, v := range routes {
				routes[i] = router.RouteReplaceDefault(v, proj)
			}
			routesJSON, err := json.MarshalIndent(def.RoutesToMap(routes), "", "  ")
			handleError(err)
			output.WriteStdout(string(routesJSON) + "\n")
			return
		}
		// table out
		data := make([][]string, 0)
		for _, route := range proj.Routes {
			route = router.RouteReplaceDefault(route, proj)
			to := route.To
			if route.Type == "upstream" {
				to = route.Upstream
			}
			data = append(data, []string{
				route.Path,
				route.Type,
				to,
			})
		}
		drawTable(
			[]string{"Path", "Type", "Upstream / To"},
			data,
		)
	},
}

func init() {
	projectStartCmd.Flags().Bool("rebuild", false, "force rebuild of app containers")
	projectStartCmd.Flags().Bool("no-build", false, "skip building project")
	projectStartCmd.Flags().Bool("no-router", false, "skip adding routes to router")
	projectStartCmd.Flags().Bool("no-commit", false, "don't commit the container after being built")
	projectStartCmd.Flags().Bool("no-validate", false, "don't validate the project config files")
	projectStartCmd.Flags().IntP("slot", "s", 0, "set volume slot")
	projectStatusCmd.Flags().Bool("json", false, "JSON output")
	projectLogsCmd.Flags().BoolP("follow", "f", false, "follow logs")
	projectRestartCmd.Flags().Bool("rebuild", false, "force rebuild of app containers")
	projectRestartCmd.Flags().Bool("no-build", false, "skip building project")
	projectRestartCmd.Flags().Bool("no-router", false, "skip adding routes to router")
	projectRestartCmd.Flags().Bool("no-commit", false, "don't commit the container after being built")
	projectRestartCmd.Flags().Bool("no-validate", false, "don't validate the project config files")
	projectRoutesCmd.Flags().Bool("json", false, "JSON output")
	projectCmd.AddCommand(projectStartCmd)
	projectCmd.AddCommand(projectStopCmd)
	projectCmd.AddCommand(projectRestartCmd)
	projectCmd.AddCommand(projectPullCmd)
	projectCmd.AddCommand(projectDeployCmd)
	projectCmd.AddCommand(projectPostDeployCmd)
	projectCmd.AddCommand(projectPurgeCmd)
	projectCmd.AddCommand(projectConfigJSONCmd)
	projectCmd.AddCommand(projectStatusCmd)
	projectCmd.AddCommand(projectLogsCmd)
	projectCmd.AddCommand(projectRoutesCmd)
	RootCmd.AddCommand(projectCmd)
}
