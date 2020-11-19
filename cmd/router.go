package cmd

import (
	"encoding/json"
	"fmt"

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
	Use:   "list",
	Short: "List routes for project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		routesJSON, err := json.MarshalIndent(def.RoutesToMap(proj.Routes), "", "  ")
		handleError(err)
		fmt.Println(string(routesJSON))
	},
}

func init() {
	routerCmd.AddCommand(routerStartCmd)
	routerCmd.AddCommand(routerStopCmd)
	routerCmd.AddCommand(routerAddCmd)
	routerCmd.AddCommand(routerDelCmd)
	routerCmd.AddCommand(routerListCmd)
	rootCmd.AddCommand(routerCmd)
}
