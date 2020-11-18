package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

// getProject fetches the project at the current working directory.
func getProject(parseYaml bool) (*project.Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return project.LoadFromPath(cwd, parseYaml)
}

// getApp fetches an app definition.
func getApp(cmd *cobra.Command, proj *project.Project) (*def.App, error) {
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

// getService fetches a service definition.
func getService(cmd *cobra.Command, proj *project.Project, filterType []string) (*def.Service, error) {
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

// handleError handles an error.
func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Println("= ERROR =======================================")
	tracerr.PrintSourceColor(err)
	os.Exit(1)
}
