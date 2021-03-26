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
	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/router"
)

var routerCertificatesCmd = &cobra.Command{
	Use:     "certificates",
	Aliases: []string{"certificate", "cert", "cer", "c", "ssl"},
	Short:   "Manage router SSL certificates.",
}

var routerCertificatesClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"delete", "del"},
	Short:   "Delete all certificate files and start anew.",
	Run: func(cmd *cobra.Command, args []string) {
		handleError(router.ClearCertificates())
	},
}

var routerCertificatesDumpCACmd = &cobra.Command{
	Use:     "dump-ca",
	Aliases: []string{"dump"},
	Short:   "Dump the CA certificate key.",
	Run: func(cmd *cobra.Command, args []string) {
		dump, err := router.DumpCertificateCA()
		handleError(err)
		output.WriteStdout(string(dump))
	},
}

func init() {
	routerCertificatesCmd.AddCommand(routerCertificatesClearCmd)
	routerCertificatesCmd.AddCommand(routerCertificatesDumpCACmd)
	routerCmd.AddCommand(routerCertificatesCmd)
}
