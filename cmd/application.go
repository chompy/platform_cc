package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api"
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
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		proj, err := api.LoadProjectFromPath(cwd, true)
		if err != nil {
			return err
		}
		appName := appCmd.PersistentFlags().Lookup("name").Value.String()
		if appName == "" {
			appName = proj.Apps[0].Name
		}
		for _, app := range proj.Apps {
			if app.Name == appName {
				return proj.BuildApp(app)
			}
		}
		return fmt.Errorf("app '%s' not found", appName)
	},
}

var appShellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh"},
	Short:   "Shell in to an application.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		proj, err := api.LoadProjectFromPath(cwd, true)
		if err != nil {
			return err
		}
		appName := appCmd.PersistentFlags().Lookup("name").Value.String()
		if appName == "" {
			appName = proj.Apps[0].Name
		}
		for _, app := range proj.Apps {
			if app.Name == appName {
				return proj.ShellApp(app)
			}
		}
		return fmt.Errorf("app '%s' not found", appName)
	},
}

func init() {
	appCmd.PersistentFlags().StringP("name", "n", "", "name of application")
	appCmd.AddCommand(appBuildCmd)
	appCmd.AddCommand(appShellCmd)
	rootCmd.AddCommand(appCmd)
}
