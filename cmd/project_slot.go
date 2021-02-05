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
	"strconv"

	"github.com/spf13/cobra"
)

var projectSlotCmd = &cobra.Command{
	Use:     "slot",
	Aliases: []string{"slt", "slots"},
	Short:   "Manage project slots.",
}

var projectSlotDelete = &cobra.Command{
	Use:   "delete slot",
	Short: "Delete slot.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			handleError(fmt.Errorf("slot argument not provided"))
		}
		proj, err := getProject(false)
		handleError(err)
		slot, err := strconv.Atoi(args[0])
		handleError(err)
		proj.SetVolumeSlot(slot)
		handleError(proj.PurgeSlot())
	},
}

func init() {
	projectSlotCmd.AddCommand(projectSlotDelete)
	projectCmd.AddCommand(projectSlotCmd)
}
