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
	"io/ioutil"
	"os"
	"sort"

	"gitlab.com/contextualcode/platform_cc/pkg/config"
	"gitlab.com/contextualcode/platform_cc/pkg/def"
	"gitlab.com/contextualcode/platform_cc/pkg/output"

	"github.com/spf13/cobra"
)

var varCmd = &cobra.Command{
	Use:     "variable [-g global]",
	Aliases: []string{"var"},
	Short:   "Manage project variables.",
}

var varSetCmd = &cobra.Command{
	Use:     "create name value",
	Aliases: []string{"set", "update"},
	Short:   "Create/set a variable.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			handleError(fmt.Errorf("missing variable name"))
		}
		// fetch value
		fi, _ := os.Stdin.Stat()
		hasStdin := fi.Mode()&os.ModeDevice == 0
		value := ""
		if len(args) >= 2 {
			value = args[1]
		} else if hasStdin {
			valueBytes, err := ioutil.ReadAll(os.Stdin)
			handleError(err)
			value = string(valueBytes)
		}
		// use global variable
		if checkFlag(varCmd, "global") {
			gc, err := config.Load()
			handleError(err)
			if value == "" {
				// delete empty var
				gc.Variables.Delete(args[0])
			} else {
				handleError(gc.Variables.Set(args[0], value))
			}
			output.Info(fmt.Sprintf("Set global variable '%s.'", args[0]))
			handleError(config.Save(gc))
			return
		}
		// use project variable
		proj, err := getProject(false)
		handleError(err)
		if value == "" {
			// delete empty var
			handleError(proj.VarDelete(args[0]))
		} else {
			handleError(proj.VarSet(args[0], value))
		}
		handleError(proj.Save())
	},
}

var varGetCmd = &cobra.Command{
	Use:     "get name",
	Aliases: []string{"g"},
	Short:   "Get a variable.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		if len(args) == 0 {
			handleError(fmt.Errorf("missing variable key"))
		}
		// use global variable only
		if checkFlag(varCmd, "global") {
			gc, err := config.Load()
			handleError(err)
			output.WriteStdout(gc.Variables.GetString(args[0]) + "\n")
			return
		}
		// use project variable (global if project var not set)
		proj, err := getProject(false)
		handleError(err)
		out, err := proj.VarGet(args[0])
		handleError(err)
		if out == "" {
			gc, err := config.Load()
			handleError(err)
			out = gc.Variables.GetString(args[0])
		}
		output.WriteStdout(out + "\n")
	},
}

var varDelCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "d"},
	Short:   "Delete a variable.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			handleError(fmt.Errorf("missing variable key"))
		}
		// use global variable only
		if checkFlag(varCmd, "global") {
			gc, err := config.Load()
			handleError(err)
			output.Info(fmt.Sprintf("Delete global variable '%s.'", args[0]))
			gc.Variables.Delete(args[0])
			handleError(config.Save(gc))
			return
		}
		// use project var
		proj, err := getProject(false)
		handleError(err)
		handleError(proj.VarDelete(args[0]))
		handleError(proj.Save())
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
		gc, err := config.Load()
		handleError(err)
		// json out
		if checkFlag(cmd, "json") {
			// build var list
			varList := make(def.Variables)
			varList.Merge(gc.Variables)
			varList.Merge(proj.Variables)
			out, err := json.MarshalIndent(varList, "", "  ")
			handleError(err)
			output.WriteStdout(string(out) + "\n")
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
	varCmd.PersistentFlags().BoolP("global", "g", false, "Use global variable.")
	varListCmd.Flags().Bool("json", false, "JSON output")
	varCmd.AddCommand(varSetCmd)
	varCmd.AddCommand(varGetCmd)
	varCmd.AddCommand(varDelCmd)
	varCmd.AddCommand(varListCmd)
	RootCmd.AddCommand(varCmd)
}
