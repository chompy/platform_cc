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
	"sort"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

var projectFlagCmd = &cobra.Command{
	Use:     "flag",
	Aliases: []string{"flags"},
	Short:   "Manage project flags.",
}

var projectFlagListCmd = &cobra.Command{
	Use:   "list [--json]",
	Short: "List project flags.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		descs := proj.Flags.Descriptions()
		// sort flag names
		sortKeys := make([]string, 0, len(descs))
		for k := range descs {
			sortKeys = append(sortKeys, k)
		}
		sort.Strings(sortKeys)
		// json out
		if checkFlag(cmd, "json") {
			data := make(map[string]map[string]interface{})
			for _, name := range sortKeys {
				data[name] = map[string]interface{}{
					"description": descs[name],
					"enabled":     proj.HasFlag(name),
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
		for _, name := range sortKeys {
			statStr := "false"
			if proj.HasFlag(name) {
				statStr = "true"
			}
			data = append(data, []string{
				name, descs[name], statStr,
			})
		}
		drawTable(
			[]string{"Name", "Description", "Enabled"},
			data,
		)
	},
}

var projectFlagSetCmd = &cobra.Command{
	Use:   "set flag [--off] [--delete]",
	Short: "Set project flag.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		if len(args) == 0 {
			handleError(fmt.Errorf("missing flag argument"))
		}
		if !proj.Flags.IsValidFlag(args[0]) {
			handleError(fmt.Errorf("%s is not a valid flag", args[0]))
		}
		value := project.FlagOn
		if checkFlag(cmd, "delete") {
			value = project.FlagUnset
		} else if checkFlag(cmd, "off") {
			value = project.FlagOff
		}
		proj.Flags.Set(args[0], value)
		handleError(proj.Save())
	},
}

var projectFlagDelCmd = &cobra.Command{
	Use:     "unset flag",
	Aliases: []string{"delete", "del", "remove", "rm", "unset"},
	Short:   "Unset project flag, allowing global flag to declare value.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		if len(args) == 0 {
			handleError(fmt.Errorf("missing flag argument"))
		}
		if !proj.Flags.IsValidFlag(args[0]) {
			handleError(fmt.Errorf("%s is not a valid flag", args[0]))
		}
		proj.Flags.Set(args[0], project.FlagUnset)
		handleError(proj.Save())
	},
}

func init() {
	projectFlagListCmd.Flags().Bool("json", false, "JSON output")
	projectFlagCmd.AddCommand(projectFlagListCmd)
	projectFlagSetCmd.Flags().Bool("off", false, "Explictly turns flag off, override global flags")
	projectFlagSetCmd.Flags().Bool("delete", false, "Unset flag, allowing global flag to declare value")
	projectFlagCmd.AddCommand(projectFlagSetCmd)
	projectFlagCmd.AddCommand(projectFlagDelCmd)
	projectCmd.AddCommand(projectFlagCmd)
}
