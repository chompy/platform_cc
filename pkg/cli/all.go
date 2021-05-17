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
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/contextualcode/platform_cc/pkg/output"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/pkg/container"
	"gitlab.com/contextualcode/platform_cc/pkg/router"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Manage all projects.",
}

var allStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all Platform.CC containers.",
	Run: func(cmd *cobra.Command, args []string) {
		containerHandler, err := getContainerHandler()
		handleError(err)
		handleError(containerHandler.AllStop())
	},
}

var allPurgeCmd = &cobra.Command{
	Use:   "purge [--global]",
	Short: "Purge all Platform.CC data.",
	Run: func(cmd *cobra.Command, args []string) {
		output.WriteStderr("!!! WARNING: PURGING ALL PLATFORM.CC DATA IN 5 SECONDS !!!\n")
		time.Sleep(5 * time.Second)
		containerHandler, err := getContainerHandler()
		handleError(err)
		handleError(containerHandler.AllPurge(
			checkFlag(cmd, "global"),
		))
	},
}

var allStatusCmd = &cobra.Command{
	Use:     "status [--json]",
	Aliases: []string{"stat", "st"},
	Short:   "Show the status of all Platform.CC containers.",
	Run: func(cmd *cobra.Command, args []string) {
		containerHandler, err := getContainerHandler()
		handleError(err)
		// retrieve status
		stats, err := containerHandler.AllStatus()
		handleError(err)
		// json out
		if checkFlag(cmd, "json") {
			out, err := json.Marshal(stats)
			handleError(err)
			output.WriteStdout(string(out) + "\n")
			return
		}
		// table out
		data := make([][]string, 0)
		for _, s := range stats {
			ipAddrStr := "n/a"
			if s.IPAddress != "" {
				ipAddrStr = s.IPAddress
			}
			slot := "n/a"
			if s.Slot > 0 {
				slot = fmt.Sprintf("%d", s.Slot)
			}
			if s.ProjectID == "" && s.Name == "" && s.Image == router.GetContainerConfig().Image {
				s.ProjectID = "-"
				s.Name = "router"
				slot = "n/a"
				s.ObjectType = container.ObjectContainerRouter
				s.Type = "router"
			}
			serviceType := s.Type
			if s.Committed {
				serviceType = "[c] " + s.Type
			}
			data = append(data, []string{
				s.ProjectID,
				fmt.Sprintf("[%s] %s", string(s.ObjectType), s.Name),
				serviceType,
				s.State,
				slot,
				ipAddrStr,
			})
		}
		drawTable(
			[]string{"Project ID", "Name", "Type", "State", "Slot", "IP Address"},
			data,
		)
		drawKeys()
		output.WriteStdout("\n")
	},
}

func init() {
	allCmd.AddCommand(allStopCmd)
	allPurgeCmd.Flags().Bool("global", false, "JSON output")
	allCmd.AddCommand(allPurgeCmd)
	allStatusCmd.Flags().Bool("json", false, "JSON output")
	allCmd.AddCommand(allStatusCmd)
	RootCmd.AddCommand(allCmd)
}
