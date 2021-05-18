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
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
)

const shareLogUpload = "https://pastebin.chompy.me/documents"

const shareLogDownload = "https://pastebin.chompy.me"

var shareLogCmd = &cobra.Command{
	Use:     "sharelog",
	Aliases: []string{"share-log"},
	Short:   "Upload current project log and provide shareable URL.",
	Run: func(cmd *cobra.Command, args []string) {
		logFile, err := os.Open(output.LogPath)
		handleError(err)
		resp, err := http.Post(
			shareLogUpload,
			"text/plain",
			logFile,
		)
		handleError(err)
		defer resp.Body.Close()
		resRaw, err := ioutil.ReadAll(resp.Body)
		handleError(err)
		res := map[string]interface{}{}
		handleError(json.Unmarshal(resRaw, &res))
		if res["key"] == nil {
			handleError(fmt.Errorf("an unknown error occured while uploading logs"))
		}
		output.WriteStdout(
			shareLogDownload + "/" + res["key"].(string) + "\n",
		)
	},
}

func init() {
	RootCmd.AddCommand(shareLogCmd)
}
