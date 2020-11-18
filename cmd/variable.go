package cmd

import (
	"encoding/json"
	"fmt"

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
	rootCmd.AddCommand(varCmd)
}
