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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitlab.com/contextualcode/platform_cc/pkg/output"

	"github.com/spf13/cobra"
)

const releasesURL = "https://platform.cc/releases/"
const latestVersionURL = releasesURL + "version"

var updateCmd = &cobra.Command{
	Use:     "update [-d dev]",
	Aliases: []string{"self-update", "upgrade"},
	Short:   "Update Platform.cc to the latest version.",
	Run: func(cmd *cobra.Command, args []string) {
		done := output.Duration("Upgrade to latest version.")
		// check latest version
		versionNo := "dev"
		if !checkFlag(cmd, "dev") {
			versionNo = fetchLatestVersionNumber()
		}
		if versionNo == "" {
			handleError(fmt.Errorf("could not determine latest version"))
		}
		// ensure we don't already have latest version
		if strings.ToLower(RootCmd.Version) == versionNo && versionNo != "dev" {
			output.Info(fmt.Sprintf("Current version (%s) is already the latest version.", strings.ToUpper(versionNo)))
			done()
			return
		}
		downloadVersion(versionNo)
		done()
	},
}

func fetchLatestVersionNumber() string {
	resp, err := http.Get(latestVersionURL)
	handleError(err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	return strings.TrimSpace(string(data))
}

func versionURL(version string) string {
	return releasesURL + strings.ToLower(version) + "/" + runtime.GOOS + "_" + runtime.GOARCH
}

func downloadVersion(version string) {
	done := output.Duration(fmt.Sprintf("Downloading version %s.", version))
	execPath, err := os.Executable()
	handleError(err)
	url := versionURL(version)
	resp, err := http.Get(url)
	handleError(err)
	defer resp.Body.Close()
	tmpPath := filepath.Join(os.TempDir(), "_pcc")
	execFile, err := os.Create(tmpPath)
	handleError(err)
	defer execFile.Close()
	_, err = io.Copy(execFile, resp.Body)
	handleError(err)
	handleError(os.Rename(tmpPath, execPath))
	handleError(os.Chmod(execPath, 0755))
	done()
}

func init() {
	updateCmd.Flags().BoolP("dev", "d", false, "Fetch latest development release.")
	RootCmd.AddCommand(updateCmd)
}
