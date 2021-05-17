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

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/pkg/output"
)

var projectValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a project.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)

		errs := proj.Validate()
		if len(errs) > 0 {
			output.Warn(fmt.Sprintf("%d validation error(s) found.", len(errs)))
			output.IndentLevel++
			for _, err := range errs {
				output.Warn(err.Error())
			}
			output.IndentLevel--
		}

	},
}

func init() {
	projectCmd.AddCommand(projectValidateCmd)
}
