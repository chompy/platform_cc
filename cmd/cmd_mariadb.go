package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api"
)

var mariadbCmd = &cobra.Command{
	Use:     "mariadb [-n name]",
	Aliases: []string{"mysql", "db"},
	Short:   "Subcommand for managing MariaDB/MySQL.",
}

var mariadbDumpCmd = &cobra.Command{
	Use:   "dump database",
	Short: "Make a database dump.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing database argument")
		}
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		service, err := getService(mariadbCmd, proj, api.MariadbService{})
		if err != nil {
			return err
		}
		return proj.ShellService(
			service,
			[]string{
				"sh", "-c", fmt.Sprintf("mysqldump -h 127.0.0.1 %s", args[0]),
			},
		)
	},
}

var mariadbShellCmd = &cobra.Command{
	Use:     "shell [-d database]",
	Aliases: []string{"sh"},
	Short:   "Shell in to database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proj, err := getProject(true)
		if err != nil {
			return err
		}
		service, err := getService(mariadbCmd, proj, api.MariadbService{})
		if err != nil {
			return err
		}
		db := cmd.PersistentFlags().Lookup("database").Value.String()
		shCmd := make([]string, 2)
		shCmd[0] = "mysql"
		shCmd[1] = "-h127.0.0.1"
		if db != "" {
			shCmd = append(shCmd, fmt.Sprintf("-D%s", db))
		}
		return proj.ShellService(
			service,
			shCmd,
		)
	},
}

func init() {
	mariadbShellCmd.PersistentFlags().StringP("database", "d", "", "name of database")
	mariadbCmd.PersistentFlags().StringP("service", "s", "", "name of service")
	mariadbCmd.AddCommand(mariadbDumpCmd)
	mariadbCmd.AddCommand(mariadbShellCmd)
	rootCmd.AddCommand(mariadbCmd)
}
