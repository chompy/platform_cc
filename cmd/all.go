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
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/docker"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Manage all projects.",
}

var allStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all PCC containers.",
	Run: func(cmd *cobra.Command, args []string) {
		docker, err := docker.NewClient()
		if err != nil {
			handleError(err)
		}
		handleError(docker.DeleteAllContainers())
	},
}

var allPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all PCC data.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("!!! WARNING: PURGING ALL PLATFORM.CC DATA IN 5 SECONDS !!!")
		time.Sleep(5 * time.Second)
		docker, err := docker.NewClient()
		handleError(err)
		handleError(docker.DeleteAllContainers())
		handleError(docker.DeleteAllVolumes())
		handleError(docker.DeleteNetwork())
	},
}

func init() {
	allCmd.AddCommand(allStopCmd)
	allCmd.AddCommand(allPurgeCmd)
	RootCmd.AddCommand(allCmd)
}
