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
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gitlab.com/contextualcode/platform_cc/api/project"
	"gitlab.com/contextualcode/platform_cc/api/router"
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

// getDef returns the best matching definition.
func getDef(name string, proj *project.Project) (interface{}, error) {
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

// getDefFromCommand returns the definition for the current command.
func getDefFromCommand(cmd *cobra.Command, proj *project.Project) (interface{}, error) {
	name := cmd.PersistentFlags().Lookup("service").Value.String()
	if name == "" {
		name = proj.Apps[0].Name
	}
	return getDef(name, proj)
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
	if name == "" {
		return def.Service{}, fmt.Errorf("no %s service was found", strings.Join(filterType, ","))
	}
	return def.Service{}, fmt.Errorf("service '%s' not found", name)
}

// handleError handles an error.
func handleError(err error) {
	output.Verbose = checkFlag(RootCmd, "verbose")
	output.Error(err)
}

// drawTable draws an ASCII table to stdout.
func drawTable(head []string, data [][]string) {
	if len(data) == 0 {
		output.WriteStdout("=== NO DATA ===\n")
		return
	}
	w, _, _ := term.GetSize(int(os.Stdin.Fd()))
	if w == 0 {
		w = 256
	}
	w -= (len(head) * 4)
	truncateString := func(size int, value string) string {
		if len(value) <= size {
			return value
		}
		if size <= 4 {
			return string(value[0]) + "..."
		}
		return value[0:size-3] + "..."
	}
	// calculate column widths
	oColWidths := make([]int, len(data[0]))
	for i, sv := range head {
		if len(sv) > oColWidths[i] {
			oColWidths[i] = len(sv)
		}
	}
	for _, v := range data {
		for i, sv := range v {
			if len(sv) > oColWidths[i] {
				oColWidths[i] = len(sv)
			}
		}
	}
	// calculate max display width per colume
	mColWidths := make([]int, len(oColWidths))
	for i := range mColWidths {
		mColWidths[i] = w / len(oColWidths)
	}
	for i := range mColWidths {
		ni := i + 1
		if ni >= len(mColWidths) {
			ni = 0
		}
		if oColWidths[i] < mColWidths[i] {
			diff := mColWidths[i] - oColWidths[i]
			mColWidths[i] -= diff
			mColWidths[ni] += diff
		}
	}
	// truncate head cols
	for i := range head {
		head[i] = truncateString(mColWidths[i], head[i])
	}
	// truncate data cols
	for i := range data {
		for j := range data[i] {
			data[i][j] = truncateString(mColWidths[j], data[i][j])
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(head)
	table.SetAutoWrapText(true)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	output.WriteStdout("\n")
	table.Render()
	output.WriteStdout("\n")
}

func drawKeys() {
	output.WriteStdout("[a] = application\t\t[s] = service\n")
	output.WriteStdout("[w] = worker\t\t\t[r] = router\n")
	output.WriteStdout("[c] = committed\n")
}

// getContainerHandler returns container handler.
func getContainerHandler() (container.Interface, error) {
	// TODO configurable
	return container.NewDocker()
}

// checkFlag returns true if given flag is set.
func checkFlag(cmd *cobra.Command, name string) bool {
	if cmd == nil {
		return false
	}
	flag := cmd.Flags().Lookup(name)
	return flag != nil && flag.Value.String() != "false"
}

func projectStart(cmd *cobra.Command, p *project.Project, slot int) {
	// get project
	if p == nil {
		var err error
		p, err = getProject(true)
		handleError(err)
	}
	// determine volume slot
	var err error
	if slot < 0 {
		slot, err = strconv.Atoi(cmd.Flags().Lookup("slot").Value.String())
		handleError(err)
	}
	p.SetSlot(slot)
	// set no commit
	if p.HasFlag(project.DisableAutoCommit) || checkFlag(cmd, "no-commit") {
		p.SetNoCommit()
	}
	// set no build
	if checkFlag(cmd, "no-build") {
		p.SetNoBuild()
	}
	// validate
	if !checkFlag(cmd, "no-validate") {
		valErrs := p.Validate()
		if len(valErrs) > 0 {
			output.ErrorText(fmt.Sprintf("Validation failed with %d error(s).", len(valErrs)))
			output.IndentLevel++
			for _, e := range valErrs {
				output.ErrorText(e.Error())
			}
			return
		}
	}
	// delete commits for rebuild
	if checkFlag(cmd, "rebuild") && !checkFlag(cmd, "no-build") {
		delComDone := output.Duration("Delete commits.")
		for _, app := range p.Apps {
			c := p.NewContainer(app)
			handleError(c.DeleteCommit())
		}
		delComDone()
	}
	// start project
	handleError(p.Start())
	// start router
	if !checkFlag(cmd, "no-router") {
		handleError(router.Start())
		handleError(router.AddProjectRoutes(p))
	}
}

// commandIntro displays introduction information about Platform.CC
func commandIntro(version string) {
	output.WriteStdout(output.Color(strings.Repeat("=", 32), 32) + "\n")
	output.WriteStdout(" PLATFORM.CC BY CONTEXTUAL CODE \n")
	output.WriteStdout(output.Color(strings.Repeat("=", 32), 32) + "\n")
	output.WriteStdout(output.Color(" VERSION "+version, 35) + "\n")
}
