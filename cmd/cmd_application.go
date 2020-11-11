package cmd

import (
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:     "application [-n name]",
	Aliases: []string{"app"},
	Short:   "Subcommand for managing applications.",
}

var appBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an application.",
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
	appCmd.AddCommand(appShellCmd)
	rootCmd.AddCommand(appCmd)
}
