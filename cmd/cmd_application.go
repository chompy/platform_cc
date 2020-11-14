package cmd

import (
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:     "application [-n name]",
	Aliases: []string{"app"},
	Short:   "Manage applications.",
}

var appBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run application build commands.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		app, err := getApp(appCmd, proj)
		if err != nil {
			return err
		}
		return proj.BuildApp(app)
	},
}

var appDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run application deploy hook.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		app, err := getApp(appCmd, proj)
		if err != nil {
			return err
		}
		return proj.DeployApp(app)
	},
}

var appShellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh"},
	Short:   "Shell in to an application.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		app, err := getApp(appCmd, proj)
		if err != nil {
			return err
		}
		return proj.ShellApp(app)
	},
}

func init() {
	appCmd.PersistentFlags().StringP("name", "n", "", "name of application")
	appCmd.AddCommand(appBuildCmd)
	appCmd.AddCommand(appDeployCmd)
	appCmd.AddCommand(appShellCmd)
	rootCmd.AddCommand(appCmd)
}
