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
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

var appPrefix = []string{"app-", "a-", "application-"}
var servicePrefix = []string{"ser-", "s-", "service-"}
var workerPrefix = []string{"wor-", "w-", "worker-"}

// getProject fetches the project at the current working directory.
func getProject(parseYaml bool) (*project.Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	proj, err := project.LoadFromPath(cwd, parseYaml)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return proj, err
}

func getDef(cmd *cobra.Command, proj *project.Project) (interface{}, error) {
	name := cmd.PersistentFlags().Lookup("name").Value.String()
	if name == "" {
		name = proj.Apps[0].Name
	}
	// find via app prefix
	for _, prefix := range appPrefix {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
			for _, d := range proj.Apps {
				if d.Name == name {
					return d, nil
				}
			}
			return nil, tracerr.Errorf("could not find application definition for %s", name)
		}
	}
	// find via service prefix
	for _, prefix := range servicePrefix {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
			for _, d := range proj.Services {
				if d.Name == name {
					return d, nil
				}
			}
			return nil, tracerr.Errorf("could not find service definition for %s", name)
		}
	}
	// find via worker prefix
	for _, prefix := range workerPrefix {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
			for _, d := range proj.Apps {
				for _, w := range d.Workers {
					if w.Name == name {
						return w, nil
					}
				}
			}
			return nil, tracerr.Errorf("could not find service definition for %s", name)
		}
	}
	// find from name with no prefix (order app > service > worker)
	for _, d := range proj.Apps {
		if d.Name == name {
			return d, nil
		}
	}
	for _, d := range proj.Services {
		if d.Name == name {
			return d, nil
		}
	}
	for _, d := range proj.Apps {
		for _, w := range d.Workers {
			if w.Name == name {
				return w, nil
			}
		}
	}
	return nil, tracerr.Errorf("could not find definition for %s", name)
}

// getService fetches a service definition.
func getService(cmd *cobra.Command, proj *project.Project, filterType []string) (def.Service, error) {
	name := cmd.PersistentFlags().Lookup("service").Value.String()
	for _, serv := range proj.Services {
		for _, t := range filterType {
			if (serv.Name == name || name == "") && t == serv.GetTypeName() {
				return serv, nil
			}
		}
	}
	return def.Service{}, fmt.Errorf("service '%s' not found", name)
}

// handleError handles an error.
func handleError(err error) {
	output.Error(err)
}

// drawTable draws an ASCII table to stdout.
func drawTable(head []string, data [][]string) {
	truncateString := func(size int, value string) string {
		if len(value) <= size {
			return value
		}
		return value[0:size-3] + "..."
	}
	for i := range head {
		head[i] = truncateString(32, head[i])
	}
	for i := range data {
		for j := range data[i] {
			data[i][j] = truncateString(256/len(data[i]), data[i][j])
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(head)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	println("")
	table.Render()
	println("")
}
