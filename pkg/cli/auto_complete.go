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
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/project"
)

// listCommandAliases returns a list of every possible command combination with aliases.
func listCommandAliases(cmd *cobra.Command) []string {
	out := make([]string, 0)
	for _, child := range cmd.Commands() {
		aliases := child.Aliases
		aliases = append(aliases, child.Name())
		if !child.HasSubCommands() {
			out = append(out, aliases...)
			continue
		}
		for _, alias := range aliases {
			childChildCmds := listCommandAliases(child)
			for _, childChildCmd := range childChildCmds {
				out = append(out, alias+":"+childChildCmd)
			}
		}
	}
	return out
}

// filterListCommandAliases returns a list of every possible command combination with given filter.
func filterListCommandAliases(cmd *cobra.Command, filter string) []string {
	allCommands := listCommandAliases(cmd)
	filterArgSplit := strings.Split(filter, ":")
	out := make([]string, 0)
	for _, cmdStr := range allCommands {
		if cmdStr == filter {
			return []string{}
		} else if filter == "" || strings.HasPrefix(cmdStr, filter) {
			// if a completed parent command has already been typed out then don't
			// include it in results
			if len(filterArgSplit) > 1 {
				for _, arg := range filterArgSplit[0 : len(filterArgSplit)-1] {
					cmdStr = strings.TrimPrefix(cmdStr, arg+":")
				}
			}
			// check for duplicate before adding
			hasOut := false
			for _, v := range out {
				if v == cmdStr {
					hasOut = true
					break
				}
			}
			if !hasOut {
				out = append(out, cmdStr)
			}
		}
	}
	return out
}

// listServices returns list of all service names with all possible filter prefixes.
func listServices() []string {
	out := make([]string, 0)
	p, err := getProject(true)
	if err != nil {
		return []string{}
	}
	for _, d := range p.Services {
		out = append(out, d.Name)
		for _, prefix := range servicePrefix {
			out = append(out, prefix+d.Name)
		}
	}
	for _, d := range p.Apps {
		out = append(out, d.Name)
		for _, prefix := range appPrefix {
			out = append(out, prefix+d.Name)
		}
		if p.HasFlag(project.EnableWorkers) {
			for _, w := range d.Workers {
				out = append(out, d.Name)
				for _, prefix := range workerPrefix {
					out = append(out, prefix+w.Name)
				}
			}
		}
	}
	return out
}

func getFlagValue(cmd *cobra.Command, name string, args []string) (bool, string) {
	cmd.ParseFlags(args)
	flag := cmd.Flag(name)
	if flag == nil {
		return false, ""
	}
	if flag.Value.String() != "" {
		return true, flag.Value.String()
	}
	if len(args) > 0 && (args[len(args)-1] == "-"+flag.Shorthand || args[len(args)-1] == "--"+flag.Name) {
		return true, ""
	}
	return false, ""
}

// listServicesFilter returns list of services filtered by given value.
func listServicesFilter(name string) []string {
	allServices := listServices()
	out := make([]string, 0)
	for _, service := range allServices {
		if name == service {
			return []string{}
		} else if strings.HasPrefix(service, name) {
			out = append(out, service)
		}
	}
	return out
}

// handleServiceFlag outputs list of available services based on given command's flag.
func handleServiceFlag(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	hasFlag, fv := getFlagValue(cmd, "service", args)
	if hasFlag {
		output.WriteStdout(strings.Join(listServicesFilter(fv), " "))
	}
}

// listDatabases returns list of all mysql databases.
func listDatabases() []string {
	p, err := getProject(true)
	if err != nil {
		return []string{}
	}
	out := make([]string, 0)
	for _, d := range p.Services {
		switch d.GetTypeName() {
		case "mysql", "mariadb":
			{
				for _, name := range d.Configuration["schemas"].([]interface{}) {
					out = append(out, name.(string))
				}
				break
			}
		}
	}
	return out
}

// listDatabasesFilter returns list of databases filtered by given value.
func listDatabasesFilter(name string) []string {
	allDatabases := listDatabases()
	out := make([]string, 0)
	for _, database := range allDatabases {
		if database == name {
			return []string{}
		} else if strings.HasPrefix(database, name) {
			out = append(out, database)
		}
	}
	return out
}

// handleDatabaseFlag outputs list of databases based on database flag.
func handleDatabaseFlag(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	hasFlag, fv := getFlagValue(cmd, "database", args)
	if hasFlag {
		output.WriteStdout(strings.Join(listDatabasesFilter(fv), " "))
	}
}

// AutoCompleteListCmd list every possible command for Bash auto-complete.
var AutoCompleteListCmd = &cobra.Command{
	Hidden: true,
	Use:    "_ac",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		output.Logging = false
		// format args so that the command is only one arg
		// bash seems to treat ":" as a seperate arg
		if len(args) > 0 {
			argStr := ""
			for _, arg := range strings.Split(args[0], ":") {
				argStr += strings.TrimSpace(arg) + ":"
			}
			argStr = strings.TrimSuffix(argStr, ":")
			args = strings.Split(argStr, " ")
		}
		if len(args) <= 1 {
			// provide list of matching commands when a partial or no command is entered
			argStr := ""
			if len(args) > 0 {
				argStr = args[0]
			}
			out := filterListCommandAliases(RootCmd, argStr)
			output.WriteStdout(strings.Join(out, " "))
		} else if len(args) > 1 {
			// otherwise provide auto complete for command flags
			// locate the command entered
			findCmd, _, err := RootCmd.Find(strings.Split(args[0], ":"))
			if err != nil {
				return
			}
			// display list of services to container commands
			if findCmd.Parent() == containerCmd {
				handleServiceFlag(containerCmd, args)
				return
			} else if findCmd.Parent() == databaseCmd {
				handleServiceFlag(databaseCmd, args)
				handleDatabaseFlag(databaseCmd, args)
				return
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(AutoCompleteListCmd)
}
