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
)

const (
	databaseMySQL    int = 1
	databasePostgres int = 2
)

var databaseTypeNames = map[int][]string{
	databaseMySQL:    []string{"mysql", "mariadb"},
	databasePostgres: []string{"postgresql"},
}

func getDatabaseTypeNames() []string {
	out := make([]string, 0)
	for _, v := range databaseTypeNames {
		out = append(out, v...)
	}
	return out
}

func matchDatabaseTypeName(name string) int {
	for k, v := range databaseTypeNames {
		for _, vv := range v {
			if vv == name {
				return k
			}
		}
	}
	return 0
}

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
		service, err := getService(databaseCmd, proj, getDatabaseTypeNames())
		handleError(err)
		database := databaseCmd.PersistentFlags().Lookup("database").Value.String()
		if database == "" {
			handleError(fmt.Errorf("must provide a database to dump"))
		}
		shellCmd := "true"
		switch matchDatabaseTypeName(service.GetTypeName()) {
		case databaseMySQL:
			{
				shellCmd = fmt.Sprintf(
					"mysqldump --password=$(cat /mnt/data/.mysql-password) %s",
					database,
				)
				break
			}
		case databasePostgres:
			{
				shellCmd = fmt.Sprintf(
					"PGPASSWORD=main pg_dump -U main -h 127.0.0.1 %s",
					database,
				)
			}
		}
		c := proj.NewContainer(service)
		handleError(c.Shell(
			"root",
			[]string{
				"sh", "-c", shellCmd,
			},
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
		service, err := getService(databaseCmd, proj, getDatabaseTypeNames())
		handleError(err)
		shellCmd := "true"
		database := databaseCmd.PersistentFlags().Lookup("database").Value.String()
		switch matchDatabaseTypeName(service.GetTypeName()) {
		case databaseMySQL:
			{
				shellCmd = "mysql --password=$(cat /mnt/data/.mysql-password)"
				if database != "" {
					shellCmd += fmt.Sprintf(" -D%s", database)
				}
				break
			}
		case databasePostgres:
			{
				shellCmd = "PGPASSWORD=main psql -U main -h 127.0.0.1"
				if database != "" {
					shellCmd += fmt.Sprintf(" --dbname=\"%s\"", database)
				}
			}
		}
		c := proj.NewContainer(service)
		handleError(c.Shell(
			"root",
			[]string{
				"sh", "-c", shellCmd,
			},
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
