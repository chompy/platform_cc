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
		fmt.Println(out)
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
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List project variables.",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Enable = false
		proj, err := getProject(false)
		if err != nil {
			return err
		}
		out, err := json.MarshalIndent(proj.Variables, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	varCmd.AddCommand(varSetCmd)
	varCmd.AddCommand(varGetCmd)
	varCmd.AddCommand(varDelCmd)
	varCmd.AddCommand(varListCmd)
	RootCmd.AddCommand(varCmd)
}
