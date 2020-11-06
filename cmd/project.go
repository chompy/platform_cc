package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Subcommand for managing projects.",
}

var projectStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		proj, err := api.LoadProjectFromPath(cwd, true)
		if err != nil {
			return err
		}
		return proj.Start()
	},
}

var projectStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		proj, err := api.LoadProjectFromPath(cwd, false)
		if err != nil {
			return err
		}
		return proj.Stop()
	},
}

var projectPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		proj, err := api.LoadProjectFromPath(cwd, false)
		if err != nil {
			return err
		}
		return proj.Purge()
	},
}

var projectConfigJSONCmd = &cobra.Command{
	Use:   "configjson",
	Short: "Dump config.json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		proj, err := api.LoadProjectFromPath(cwd, true)
		if err != nil {
			return err
		}
		configJSON, err := proj.BuildConfigJSON(&proj.Apps[0])
		if err != nil {
			return err
		}
		fmt.Println(string(configJSON))
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectStartCmd)
	projectCmd.AddCommand(projectStopCmd)
	projectCmd.AddCommand(projectPurgeCmd)
	projectCmd.AddCommand(projectConfigJSONCmd)
	rootCmd.AddCommand(projectCmd)
}
