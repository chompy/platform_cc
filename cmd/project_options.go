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

	"gitlab.com/contextualcode/platform_cc/api/output"

	"gitlab.com/contextualcode/platform_cc/api/project"

	"github.com/spf13/cobra"
)

var projectOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"opt", "opts", "option"},
	Short:   "Manage project options.",
}

var projectOptionListCmd = &cobra.Command{
	Use:   "list [--json]",
	Short: "List project options.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		descs := project.ListOptionDescription()
		// json out
		if checkFlag(cmd, "json") {
			data := make(map[string]map[string]interface{})
			for opt, desc := range descs {
				data[string(opt)] = map[string]interface{}{
					"description": desc,
					"default":     opt.DefaultValue(),
					"value":       proj.GetOption(opt),
				}
			}
			out, err := json.MarshalIndent(
				data,
				"",
				"  ",
			)
			handleError(err)
			output.WriteStdout(string(out) + "\n")
			return
		}
		// table out
		data := make([][]string, 0)
		for opt, desc := range descs {
			data = append(data, []string{
				string(opt), desc, opt.DefaultValue(), proj.GetOption(opt),
			})
		}
		drawTable(
			[]string{"Name", "Description", "Default", "Value"},
			data,
		)
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
		// itterate possible options
		for _, opt := range project.ListOptions() {
			if string(opt) == args[0] {
				if proj.Options == nil {
					proj.Options = make(map[project.Option]string)
				}
				handleError(opt.Validate(args[1]))
				proj.Options[opt] = args[1]
				handleError(proj.Save())
				return
			}
		}
		handleError(fmt.Errorf("'%s' is not a valid option", args[0]))
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
		// itterate possible options
		for _, opt := range project.ListOptions() {
			if string(opt) == args[0] {
				if proj.Options == nil {
					proj.Options = make(map[project.Option]string)
				}
				proj.Options[opt] = ""
				handleError(proj.Save())
				return
			}
		}
		handleError(fmt.Errorf("'%s' is not a valid option", args[0]))
	},
}

func init() {
	projectOptionListCmd.Flags().Bool("json", false, "JSON output")
	projectOptionsCmd.AddCommand(projectOptionListCmd)
	projectOptionsCmd.AddCommand(projectOptionSetCmd)
	projectOptionsCmd.AddCommand(projectOptionDelCmd)
	projectCmd.AddCommand(projectOptionsCmd)
}
