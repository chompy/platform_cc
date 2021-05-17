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
	"strconv"
	"time"

	"gitlab.com/contextualcode/platform_cc/pkg/output"

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
		if len(args) != 1 {
			handleError(fmt.Errorf("slot argument not provided"))
		}
		proj, err := getProject(false)
		handleError(err)
		slot, err := strconv.Atoi(args[0])
		handleError(err)
		proj.SetSlot(slot)
		if slot <= 1 {
			output.WriteStderr("!!! WARNING: DELETING ALL DATA IN SLOT 1 IN 5 SECONDS !!!\n")
			time.Sleep(time.Second * 5)
		}
		handleError(proj.PurgeSlot())
	},
}

var projectSlotCopy = &cobra.Command{
	Use:   "copy source dest",
	Short: "Copy slot.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			handleError(fmt.Errorf("slot source and/or slot destination not provided"))
		}
		proj, err := getProject(false)
		handleError(err)
		source, err := strconv.Atoi(args[0])
		handleError(err)
		dest, err := strconv.Atoi(args[1])
		handleError(err)
		proj.SetSlot(source)
		handleError(proj.CopySlot(dest))
	},
}

func init() {
	projectSlotCmd.AddCommand(projectSlotDelete)
	projectSlotCmd.AddCommand(projectSlotCopy)
	projectCmd.AddCommand(projectSlotCmd)
}
