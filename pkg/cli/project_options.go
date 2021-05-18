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
	"encoding/json"
	"fmt"

	"gitlab.com/contextualcode/platform_cc/v2/pkg/config"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"

	"gitlab.com/contextualcode/platform_cc/v2/pkg/project"

	"github.com/spf13/cobra"
)

var projectOptionsCmd = &cobra.Command{
	Use:     "options [-g global]",
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
					"source":      proj.OptionSource(opt),
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
				string(opt), desc, proj.OptionSource(opt), opt.DefaultValue(), proj.GetOption(opt),
			})
		}
		drawTable(
			[]string{"Name", "Description", "Source", "Default", "Value"},
			data,
		)
	},
}

var projectOptionSetCmd = &cobra.Command{
	Use:   "set option value",
	Short: "Set project option.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			handleError(fmt.Errorf("missing option and/or value argument(s)"))
			return
		}
		// itterate possible options
		for _, opt := range project.ListOptions() {
			if string(opt) == args[0] {
				handleError(opt.Validate(args[1]))
				// global
				if checkFlag(projectOptionsCmd, "global") {
					gc, err := config.Load()
					handleError(err)
					output.Info(fmt.Sprintf("Set global option '%s.'", string(opt)))
					gc.Options[string(opt)] = args[1]
					handleError(config.Save(gc))
					return
				}
				// local
				proj, err := getProject(false)
				handleError(err)
				output.Info(fmt.Sprintf("Set option '%s.'", string(opt)))
				if proj.Options == nil {
					proj.Options = make(map[project.Option]string)
				}
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
		if len(args) == 0 {
			handleError(fmt.Errorf("missing option argument"))
			return
		}
		// itterate possible options
		for _, opt := range project.ListOptions() {
			if string(opt) == args[0] {
				// global
				if checkFlag(projectOptionsCmd, "global") {
					gc, err := config.Load()
					handleError(err)
					output.Info(fmt.Sprintf("Delete global option '%s.'", string(opt)))
					gc.Options[string(opt)] = ""
					handleError(config.Save(gc))
					return
				}
				// local
				proj, err := getProject(false)
				handleError(err)
				output.Info(fmt.Sprintf("Delete option '%s.'", string(opt)))
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
	projectOptionsCmd.PersistentFlags().BoolP("global", "g", false, "Use global options.")
	projectOptionListCmd.Flags().Bool("json", false, "JSON output")
	projectOptionsCmd.AddCommand(projectOptionListCmd)
	projectOptionsCmd.AddCommand(projectOptionSetCmd)
	projectOptionsCmd.AddCommand(projectOptionDelCmd)
	projectCmd.AddCommand(projectOptionsCmd)
}
