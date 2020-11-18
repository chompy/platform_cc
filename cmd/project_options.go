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
