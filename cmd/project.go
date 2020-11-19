package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/router"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj", "p"},
	Short:   "Manage projects.",
}

var projectStartCmd = &cobra.Command{
	Use:   "start [--no-build] [--no-router]",
	Short: "Start a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.Start())
		noBuildFlag := cmd.Flags().Lookup("no-build")
		if noBuildFlag == nil || noBuildFlag.Value.String() == "false" {
			handleError(proj.Build())
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

var projectBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		handleError(proj.Build())
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
		proj, err := getProject(true)
		handleError(err)
		configJSON, err := proj.BuildConfigJSON(&proj.Apps[0])
		handleError(err)
		fmt.Println(string(configJSON))
	},
}

func init() {
	projectStartCmd.Flags().Bool("no-build", false, "skip building project")
	projectStartCmd.Flags().Bool("no-router", false, "skip adding routes to router")
	projectCmd.AddCommand(projectStartCmd)
	projectCmd.AddCommand(projectStopCmd)
	projectCmd.AddCommand(projectRestartCmd)
	projectCmd.AddCommand(projectBuildCmd)
	projectCmd.AddCommand(projectDeployCmd)
	projectCmd.AddCommand(projectPurgeCmd)
	projectCmd.AddCommand(projectConfigJSONCmd)
	rootCmd.AddCommand(projectCmd)
}
