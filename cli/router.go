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
	"encoding/json"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/router"
)

var routerCmd = &cobra.Command{
	Use:     "router",
	Aliases: []string{"r"},
	Short:   "Manage router.",
}

var routerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start router.",
	Run: func(cmd *cobra.Command, args []string) {
		handleError(router.Start())
	},
}

var routerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop router.",
	Run: func(cmd *cobra.Command, args []string) {
		handleError(router.Stop())
	},
}

var routerResetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"restart"},
	Short:   "Remove all routes from the router.",
	Run: func(cmd *cobra.Command, args []string) {
		handleError(router.Stop())
		handleError(router.Start())
	},
}

var routerAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add project to router.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(router.AddProjectRoutes(proj))
	},
}

var routerDelCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "remove"},
	Short:   "Delete project from router.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(router.DeleteProjectRoutes(proj))
	},
}

var routerListCmd = &cobra.Command{
	Use:   "list [--json]",
	Short: "List all active routes.",
	Run: func(cmd *cobra.Command, args []string) {
		routes, err := router.ListActiveRoutes()
		handleError(err)
		// json out
		if checkFlag(cmd, "json") {
			routesJSON, err := json.MarshalIndent(def.RoutesToMap(routes), "", "  ")
			handleError(err)
			output.WriteStdout(string(routesJSON) + "\n")
			return
		}
		// table out
		data := make([][]string, 0)
		for _, route := range routes {
			to := route.To
			if route.Type == "upstream" {
				to = route.Upstream
			}
			data = append(data, []string{
				route.Attributes["project_id"],
				route.Path,
				route.Type,
				to,
			})
		}
		drawTable(
			[]string{"Project ID", "Path", "Type", "Upstream / To"},
			data,
		)
	},
}

func init() {
	routerListCmd.Flags().Bool("json", false, "JSON output")
	routerCmd.AddCommand(routerStartCmd)
	routerCmd.AddCommand(routerStopCmd)
	routerCmd.AddCommand(routerResetCmd)
	routerCmd.AddCommand(routerAddCmd)
	routerCmd.AddCommand(routerDelCmd)
	routerCmd.AddCommand(routerListCmd)
	RootCmd.AddCommand(routerCmd)
}
