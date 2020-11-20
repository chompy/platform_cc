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
				"sh", "-c",
				fmt.Sprintf(
					"mysqldump -p$(cat /mnt/data/.mysql-password) %s",
					args[0],
				),
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
		dbCmdStr := ""
		db := cmd.PersistentFlags().Lookup("database").Value.String()
		if db != "" {
			dbCmdStr = fmt.Sprintf(" -D%s", db)
		}
		shCmd := []string{
			"sh", "-c",
			fmt.Sprintf(
				"mysql -p$(cat /mnt/data/.mysql-password)%s",
				dbCmdStr,
			),
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
	RootCmd.AddCommand(mariadbCmd)
}
