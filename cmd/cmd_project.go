package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj", "p"},
	Short:   "Manage projects.",
}

var projectStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
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
		proj, err := getProject(false)
		if err != nil {
			return err
		}
		return proj.Stop()
	},
}

var projectRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		if err := proj.Stop(); err != nil {
			return err
		}
		return proj.Start()
	},
}

var projectBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		return proj.Build()
	},
}

var projectDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run deploy hooks for project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		return proj.Deploy()
	},
}

var projectPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(false)
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
		proj, err := getProject(true)
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
	projectCmd.AddCommand(projectRestartCmd)
	projectCmd.AddCommand(projectBuildCmd)
	projectCmd.AddCommand(projectDeployCmd)
	projectCmd.AddCommand(projectPurgeCmd)
	projectCmd.AddCommand(projectConfigJSONCmd)
	rootCmd.AddCommand(projectCmd)
}
