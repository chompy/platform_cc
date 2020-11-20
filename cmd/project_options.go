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
	"encoding/json"
	"fmt"

	"gitlab.com/contextualcode/platform_cc/api/project"

	"github.com/spf13/cobra"
)

var projectOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"opt", "opts"},
	Short:   "Manage project options.",
}

var projectOptionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List project options.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		descs := project.ListOptionDescription()
		data := make(map[string]map[string]interface{})
		for opt, desc := range descs {
			data[string(opt)] = map[string]interface{}{
				"description": desc,
				"default":     opt.DefaultValue(),
				"value":       opt.Value(proj.Options),
			}
		}
		out, err := json.MarshalIndent(
			data,
			"",
			"  ",
		)
		handleError(err)
		fmt.Println(string(out))
	},
}

var projectOptionSetCmd = &cobra.Command{
	Use:   "set option value",
	Short: "Set project option.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		if len(args) < 2 {
			handleError(fmt.Errorf("missing option and/or value argument(s)"))
			return
		}
		if proj.Options == nil {
			proj.Options = make(map[project.Option]string)
		}
		proj.Options[project.Option(args[0])] = args[1]
		handleError(proj.Save())
	},
}

var projectOptionDelCmd = &cobra.Command{
	Use:     "reset option",
	Aliases: []string{"del", "delete", "remove"},
	Short:   "Reset project option to default.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		if len(args) == 0 {
			handleError(fmt.Errorf("missing option argument"))
			return
		}
		if proj.Options == nil {
			proj.Options = make(map[project.Option]string)
		}
		proj.Options[project.Option(args[0])] = ""
		handleError(proj.Save())
	},
}

func init() {
	projectOptionsCmd.AddCommand(projectOptionListCmd)
	projectOptionsCmd.AddCommand(projectOptionSetCmd)
	projectOptionsCmd.AddCommand(projectOptionDelCmd)
	projectCmd.AddCommand(projectOptionsCmd)
}
