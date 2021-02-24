/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

var projectValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		count := 0
		defs := make([]interface{}, 0)
		for _, app := range proj.Apps {
			done := output.Duration(
				fmt.Sprintf("Validate app '%s.'", app.Name),
			)
			errs := app.Validate()
			if len(errs) > 0 {
				count += len(errs)
				for _, err := range errs {
					output.ErrorText(err.Error())
				}
			}
			defs = append(defs, app)
			done()
		}
		for _, serv := range proj.Services {
			done := output.Duration(
				fmt.Sprintf("Validate service '%s.'", serv.Name),
			)
			errs := serv.Validate()
			if len(errs) > 0 {
				count += len(errs)
				for _, err := range errs {
					output.Info(err.Error())
				}
			}
			defs = append(defs, serv)
			done()
		}
		for _, route := range proj.Routes {
			done := output.Duration(
				fmt.Sprintf("Validate route '%s.'", route.Path),
			)
			errs := route.Validate()
			if len(errs) > 0 {
				count += len(errs)
				for _, err := range errs {
					output.Info(err.Error())
				}
			}
			done()
		}
		if len(defs) > 0 {
			done := output.Duration("Check for invalid relationships.")
			_, err := proj.GetDefinitionStartOrder(defs)
			if err != nil {
				count++
				output.Info(err.Error())
			}
			done()
		}

		output.Info(fmt.Sprintf("Validation completed with %d errors(s).", count))
	},
}

func init() {
	projectCmd.AddCommand(projectValidateCmd)
}
