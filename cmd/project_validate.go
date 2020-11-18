package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

var projectValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a project.",
	Run: func(cmd *cobra.Command, args []string) {
		//var buf bytes.Buffer
		//log.SetOutput(&buf)
		count := 0
		fmt.Println("* Load project.")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("--> ERROR: could not determine current working directory,", err)
			return
		}
		proj, err := project.LoadFromPath(cwd, true)
		if err != nil {
			fmt.Println("--> ERROR: could not load project,", err)
			return
		}
		for _, app := range proj.Apps {
			fmt.Printf("* Validate app '%s.'\n", app.Name)
			errs := app.Validate()
			if len(errs) > 0 {
				count += len(errs)
				for _, err := range errs {
					fmt.Printf("\t- %s\n", err)
				}
			}
		}
		for _, serv := range proj.Services {
			fmt.Printf("* Validate service '%s.'\n", serv.Name)
			errs := serv.Validate()
			if len(errs) > 0 {
				count += len(errs)
				for _, err := range errs {
					fmt.Printf("\t- %s\n", err)
				}
			}
		}
		for _, route := range proj.Routes {
			fmt.Printf("* Validate route '%s.'\n", route.Path)
			errs := route.Validate()
			if len(errs) > 0 {
				count += len(errs)
				for _, err := range errs {
					fmt.Printf("\t- %s\n", err)
				}
			}
		}
		fmt.Printf("* Validation completed (%d errors(s)).\n", count)
	},
}

func init() {
	projectCmd.AddCommand(projectValidateCmd)
}
