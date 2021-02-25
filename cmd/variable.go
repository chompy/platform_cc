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
	"sort"

	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/spf13/cobra"
)

var varCmd = &cobra.Command{
	Use:     "variable",
	Aliases: []string{"var"},
	Short:   "Manage project variables.",
}

var varSetCmd = &cobra.Command{
	Use:     "set key value",
	Aliases: []string{"s"},
	Short:   "Set a variable.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(false)
		if err != nil {
			return err
		}
		if err := proj.VarSet(args[0], args[1]); err != nil {
			return err
		}
		return proj.Save()
	},
}

var varGetCmd = &cobra.Command{
	Use:     "get key",
	Aliases: []string{"g"},
	Short:   "Get a variable.",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Enable = false
		proj, err := getProject(false)
		if err != nil {
			return err
		}
		out, err := proj.VarGet(args[0])
		if err != nil {
			return err
		}
		println(out)
		return nil
	},
}

var varDelCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "d"},
	Short:   "Delete a variable.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(false)
		if err != nil {
			return err
		}
		if err := proj.VarDelete(args[0]); err != nil {
			return err
		}
		return proj.Save()
	},
}

var varListCmd = &cobra.Command{
	Use:     "list [--json]",
	Aliases: []string{"l"},
	Short:   "List project variables.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(false)
		handleError(err)
		// global config
		gc, err := def.ParseGlobalYamlFile()
		handleError(err)
		// json out
		jsonFlag := cmd.Flags().Lookup("json")
		if jsonFlag != nil && jsonFlag.Value.String() != "false" {
			// build var list
			varList := make(def.Variables)
			varList.Merge(gc.Variables)
			varList.Merge(proj.Variables)
			out, err := json.MarshalIndent(varList, "", "  ")
			handleError(err)
			println(string(out))
			return
		}

		varSources := make(map[string][]string)
		varKeys := make([]string, 0)
		for _, k := range gc.Variables.Keys() {
			varSources[k] = []string{"global", gc.Variables.GetString(k)}
		}
		for _, k := range proj.Variables.Keys() {
			varSources[k] = []string{"project", proj.Variables.GetString(k)}
		}
		for k := range varSources {
			varKeys = append(varKeys, k)
		}
		sort.Strings(varKeys)

		// table out
		data := make([][]string, 0)
		for _, k := range varKeys {
			data = append(data, []string{k, varSources[k][0], varSources[k][1]})
		}
		drawTable(
			[]string{"Key", "Source", "Value"},
			data,
		)
	},
}

func init() {
	varListCmd.Flags().Bool("json", false, "JSON output")
	varCmd.AddCommand(varSetCmd)
	varCmd.AddCommand(varGetCmd)
	varCmd.AddCommand(varDelCmd)
	varCmd.AddCommand(varListCmd)
	RootCmd.AddCommand(varCmd)
}
