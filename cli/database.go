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
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

var databaseCmd = &cobra.Command{
	Use:     "database [-s service] [-d database]",
	Aliases: []string{"db", "mysql", "mariadb"},
	Short:   "Manage database.",
}

var databaseDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Make a database dump.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		service, err := getService(databaseCmd, proj, project.GetDatabaseTypeNames())
		handleError(err)
		database := databaseCmd.PersistentFlags().Lookup("database").Value.String()
		if database == "" {
			handleError(fmt.Errorf("must provide a database to dump"))
		}
		c := proj.NewContainer(service)
		handleError(c.Shell(
			"root",
			proj.GetDatabaseDumpCommand(service, database),
		))
	},
}

var databaseShellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh", "sql"},
	Short:   "Shell in to database.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		service, err := getService(databaseCmd, proj, project.GetDatabaseTypeNames())
		handleError(err)
		database := databaseCmd.PersistentFlags().Lookup("database").Value.String()
		c := proj.NewContainer(service)
		handleError(c.Shell(
			"root",
			proj.GetDatabaseShellCommand(service, database),
		))
	},
}

func init() {
	databaseCmd.PersistentFlags().StringP("database", "d", "", "name of database")
	databaseCmd.PersistentFlags().StringP("service", "s", "", "name of service")
	databaseCmd.AddCommand(databaseDumpCmd)
	databaseCmd.AddCommand(databaseShellCmd)
	RootCmd.AddCommand(databaseCmd)
}
