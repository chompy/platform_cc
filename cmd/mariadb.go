package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mariadbTypeNames = []string{"mysql", "mariadb"}

var mariadbCmd = &cobra.Command{
	Use:     "mariadb [-n name]",
	Aliases: []string{"mysql", "db"},
	Short:   "Manage MariaDB/MySQL.",
}

var mariadbDumpCmd = &cobra.Command{
	Use:   "dump database",
	Short: "Make a database dump.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			handleError(fmt.Errorf("missing database argument"))
		}
		proj, err := getProject(true)
		handleError(err)
		service, err := getService(mariadbCmd, proj, mariadbTypeNames)
		handleError(err)
		handleError(proj.ShellService(
			service,
			[]string{
				"sh", "-c", fmt.Sprintf("mysqldump -h 127.0.0.1 %s", args[0]),
			},
		))
	},
}

var mariadbShellCmd = &cobra.Command{
	Use:     "shell [-d database]",
	Aliases: []string{"sh", "sql"},
	Short:   "Shell in to database.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		service, err := getService(mariadbCmd, proj, mariadbTypeNames)
		handleError(err)
		db := cmd.PersistentFlags().Lookup("database").Value.String()
		shCmd := make([]string, 2)
		shCmd[0] = "mysql"
		shCmd[1] = "-h127.0.0.1"
		if db != "" {
			shCmd = append(shCmd, fmt.Sprintf("-D%s", db))
		}
		handleError(proj.ShellService(
			service,
			shCmd,
		))
	},
}

func init() {
	mariadbShellCmd.PersistentFlags().StringP("database", "d", "", "name of database")
	mariadbCmd.PersistentFlags().StringP("service", "s", "", "name of service")
	mariadbCmd.AddCommand(mariadbDumpCmd)
	mariadbCmd.AddCommand(mariadbShellCmd)
	rootCmd.AddCommand(mariadbCmd)
}
