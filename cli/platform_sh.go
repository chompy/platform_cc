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
	"log"

	"github.com/spf13/cobra"
)

var platformShCmd = &cobra.Command{
	Use:     "platform-sh",
	Aliases: []string{"platform_sh", "platformsh", "psh"},
	Short:   "Manage Platform.sh options.",
}

var platformShTestCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(false)
		handleError(err)
		psh, err := proj.GetPlatformSHProject()
		handleError(err)
		log.Println(psh.ID)
	},
}

func init() {
	platformShCmd.AddCommand(platformShTestCmd)
	RootCmd.AddCommand(platformShCmd)
}
