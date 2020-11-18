package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var projectFlagCmd = &cobra.Command{
	Use:     "flag",
	Aliases: []string{"flags"},
	Short:   "Manage project flags.",
}

var projectFlagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List project flags.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		descs := proj.Flags.Descriptions()
		flags := proj.Flags.List()
		data := make(map[string]map[string]interface{})
		for name, desc := range descs {
			data[name] = map[string]interface{}{
				"description": desc,
				"enabled":     proj.Flags.Has(flags[name]),
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

var projectFlagSetCmd = &cobra.Command{
	Use:   "set flag",
	Short: "Set project flags.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		flags := proj.Flags.List()
		if len(args) == 0 {
			handleError(fmt.Errorf("missing flag argument"))
		}
		if flags[args[0]] == 0 {
			handleError(fmt.Errorf("%s is not a valid flag", args[0]))
		}
		if !proj.Flags.Has(flags[args[0]]) {
			proj.Flags.Add(flags[args[0]])
			handleError(proj.Save())
		}
	},
}

var projectFlagDelCmd = &cobra.Command{
	Use:     "remove flag",
	Aliases: []string{"delete", "del"},
	Short:   "Remove project flags.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		flags := proj.Flags.List()
		if len(args) == 0 {
			handleError(fmt.Errorf("missing flag argument"))
		}
		if flags[args[0]] == 0 {
			handleError(fmt.Errorf("%s is not a valid flag", args[0]))
		}
		if proj.Flags.Has(flags[args[0]]) {
			proj.Flags.Remove(flags[args[0]])
			handleError(proj.Save())
		}
	},
}

func init() {
	projectFlagCmd.AddCommand(projectFlagListCmd)
	projectFlagCmd.AddCommand(projectFlagSetCmd)
	projectFlagCmd.AddCommand(projectFlagDelCmd)
	projectCmd.AddCommand(projectFlagCmd)
}
