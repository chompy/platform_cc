package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api"
)

// getProject - fetch project for commands
func getProject(parseYaml bool) (*api.Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return api.LoadProjectFromPath(cwd, parseYaml)
}

// getApp - fetch app def
func getApp(cmd *cobra.Command, proj *api.Project) (*api.AppDef, error) {
	name := cmd.PersistentFlags().Lookup("name").Value.String()
	if name == "" {
		name = proj.Apps[0].Name
	}
	for _, app := range proj.Apps {
		if app.Name == name {
			return app, nil
		}
	}
	return nil, fmt.Errorf("app '%s' not found", name)
}

// getService - fetch service def
func getService(cmd *cobra.Command, proj *api.Project, filterType []string) (*api.ServiceDef, error) {
	name := cmd.PersistentFlags().Lookup("service").Value.String()
	for _, serv := range proj.Services {
		for _, t := range filterType {
			if (serv.Name == name || name == "") && t == serv.GetTypeName() {
				return serv, nil
			}
		}
	}
	return nil, fmt.Errorf("service '%s' not found", name)
}
