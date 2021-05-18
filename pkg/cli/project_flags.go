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

	"gitlab.com/contextualcode/platform_cc/v2/pkg/config"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/project"
)

var projectFlagCmd = &cobra.Command{
	Use:     "flag [-g global]",
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
					"source":      proj.FlagSource(name),
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
				name, descs[name], proj.FlagSource(name), statStr,
			})
		}
		drawTable(
			[]string{"Name", "Description", "Source", "Enabled"},
			data,
		)
	},
}

var projectFlagSetCmd = &cobra.Command{
	Use:   "set flag [--off] [--delete]",
	Short: "Set project flag.",
	Run: func(cmd *cobra.Command, args []string) {
		// validate
		proj, err := getProject(false)
		handleError(err)
		if len(args) == 0 {
			handleError(fmt.Errorf("missing flag argument"))
		}
		if !proj.Flags.IsValidFlag(args[0]) {
			handleError(fmt.Errorf("%s is not a valid flag", args[0]))
		}
		// set global
		if checkFlag(projectFlagCmd, "global") {
			gc, err := config.Load()
			handleError(err)
			// off/delete do the same thing, remove flag
			if checkFlag(cmd, "delete") || checkFlag(cmd, "off") {
				output.Info(fmt.Sprintf("Set global flag '%s' to unset.", args[0]))
				out := make([]string, 0)
				for _, flag := range gc.Flags {
					if flag == args[0] {
						continue
					}
					out = append(out, flag)
				}
				gc.Flags = out
				handleError(config.Save(gc))
				return
			}
			// check already set
			output.Info(fmt.Sprintf("Set global flag '%s' to on.", args[0]))
			for _, flag := range gc.Flags {
				if flag == args[0] {
					return
				}
			}
			gc.Flags = append(gc.Flags, args[0])
			handleError(config.Save(gc))
			return
		}
		// set local
		value := project.FlagOn
		valueStr := "on"
		if checkFlag(cmd, "delete") {
			value = project.FlagUnset
			valueStr = "unset"
		} else if checkFlag(cmd, "off") {
			value = project.FlagOff
			valueStr = "off"
		}
		output.Info(fmt.Sprintf("Set flag '%s' to %s.", args[0], valueStr))
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
		// global
		if checkFlag(projectFlagCmd, "global") {
			gc, err := config.Load()
			handleError(err)
			output.Info(fmt.Sprintf("Set global flag '%s' to unset.", args[0]))
			out := make([]string, 0)
			for _, flag := range gc.Flags {
				if flag == args[0] {
					continue
				}
				out = append(out, flag)
			}
			gc.Flags = out
			handleError(config.Save(gc))
			return
		}
		// local
		output.Info(fmt.Sprintf("Set flag '%s' to unset.", args[0]))
		proj.Flags.Set(args[0], project.FlagUnset)
		handleError(proj.Save())
	},
}

func init() {
	projectFlagCmd.PersistentFlags().BoolP("global", "g", false, "Use global flag.")
	projectFlagListCmd.Flags().Bool("json", false, "JSON output")
	projectFlagCmd.AddCommand(projectFlagListCmd)
	projectFlagSetCmd.Flags().Bool("off", false, "Explictly turns flag off, override global flags")
	projectFlagSetCmd.Flags().Bool("delete", false, "Unset flag, allowing global flag to declare value")
	projectFlagCmd.AddCommand(projectFlagSetCmd)
	projectFlagCmd.AddCommand(projectFlagDelCmd)
	projectCmd.AddCommand(projectFlagCmd)
}
