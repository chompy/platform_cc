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

package project

import (
	"fmt"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

const (
	databaseMySQL    int = 1
	databasePostgres int = 2
)

var databaseTypeNames = map[int][]string{
	databaseMySQL:    []string{"mysql", "mariadb"},
	databasePostgres: []string{"postgresql"},
}

// GetDatabaseTypeNames returns list of all service types that are considered to be a database.
func GetDatabaseTypeNames() []string {
	out := make([]string, 0)
	for _, v := range databaseTypeNames {
		out = append(out, v...)
	}
	return out
}

// MatchDatabaseTypeName matches given service type with database service type.
func MatchDatabaseTypeName(name string) int {
	for k, v := range databaseTypeNames {
		for _, vv := range v {
			if vv == name {
				return k
			}
		}
	}
	return 0
}

// GetDatabaseShellCommand returns the command to access the database shell for given definition.
func (p *Project) GetDatabaseShellCommand(d interface{}, database string) string {
	switch d.(type) {
	case def.Service:
		{
			service := d.(def.Service)
			switch MatchDatabaseTypeName(service.GetTypeName()) {
			case databaseMySQL:
				{
					shellCmd := "mysql --password=$(cat /mnt/data/.mysql-password)"
					if database != "" {
						shellCmd += fmt.Sprintf(" -D%s", database)
					}
					return shellCmd
				}
			case databasePostgres:
				{
					shellCmd := "PGPASSWORD=main psql -U main -h 127.0.0.1"
					if database != "" {
						shellCmd += fmt.Sprintf(" --dbname=\"%s\"", database)
					}
					return shellCmd
				}
			}
		}
	}
	return ""
}

// GetDatabaseDumpCommand returns the command to dump a database for given definition.
func (p *Project) GetDatabaseDumpCommand(d interface{}, database string) string {
	switch d.(type) {
	case def.Service:
		{
			service := d.(def.Service)
			switch MatchDatabaseTypeName(service.GetTypeName()) {
			case databaseMySQL:
				{
					shellCmd := fmt.Sprintf(
						"mysqldump --password=$(cat /mnt/data/.mysql-password) %s",
						database,
					)
					return shellCmd
				}
			case databasePostgres:
				{
					shellCmd := fmt.Sprintf(
						"PGPASSWORD=main pg_dump -U main -h 127.0.0.1 %s",
						database,
					)
					return shellCmd
				}
			}
		}
	}
	return ""
}

// GetPlatformSHDatabaseDumpCommand returns the command to dump a database from Platform.sh for given definition.
func (p *Project) GetPlatformSHDatabaseDumpCommand(d interface{}, database string, rels map[string]interface{}) string {
	switch d.(type) {
	case def.Service:
		{
			service := d.(def.Service)
			// find valid relationship and get creds
			rel := ""
			user := ""
			pass := ""
			host := ""
			for name, endpoint := range service.Configuration["endpoints"].(map[string]interface{}) {
				for dbName := range endpoint.(map[string]interface{})["privileges"].(map[string]interface{}) {
					if dbName == database {
						rel = name
						break
					}
				}
			}
			for _, v := range rels {
				for _, vv := range v.([]interface{}) {
					val := vv.(map[string]interface{})
					if val["service"] == service.Name && val["rel"] == rel {
						user = val["username"].(string)
						pass = val["password"].(string)
						host = val["host"].(string)
						break
					}
				}
			}
			// get dump command
			switch MatchDatabaseTypeName(service.GetTypeName()) {
			case databaseMySQL:
				{
					return fmt.Sprintf(
						`mysqldump --host="%s" -u%s --password=%s %s`,
						host,
						user,
						pass,
						database,
					)
				}
			case databasePostgres:
				{
					return fmt.Sprintf(
						"PGPASSWORD=%s pg_dump -U %s -h %s %s",
						pass,
						user,
						host,
						database,
					)
				}
			}
		}
	}
	return ""
}
