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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

var containerCmd = &cobra.Command{
	Use:     "container [-s service]",
	Aliases: []string{"cont", "c", "application", "app"},
	Short:   "Manage individual containers for applications, services, and workers.",
}

var containerAppDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run application deploy hook.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDefFromCommand(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.Deploy())
				return
			}
		}
		handleError(fmt.Errorf("can only run deploy hooks on applications"))
	},
}

var containerAppPostDeployCmd = &cobra.Command{
	Use:     "post-deploy",
	Aliases: []string{"postdeploy", "pd"},
	Short:   "Run application post-deploy hook.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDefFromCommand(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.PostDeploy())
				return
			}
		}
		handleError(fmt.Errorf("can only run post-deploy hooks on applications"))
	},
}

var containerShellCmd = &cobra.Command{
	Use:     "shell [--root] command",
	Aliases: []string{"sh"},
	Short:   "Shell in to a container.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDefFromCommand(containerCmd, proj)
		handleError(err)
		user := "root"
		switch d.(type) {
		case def.App:
			{
				user = "web"
				if cmd.PersistentFlags().Lookup("root").Value.String() == "true" {
					user = "root"
				}
				break
			}
		}
		c := proj.NewContainer(d)
		shellCmd := []string{}
		if len(args) > 0 {
			shellCmd = []string{
				"bash", "--login", "-c", strings.Join(args, " "),
			}
		}
		handleError(c.Shell(user, shellCmd))
	},
}

var containerAppCommitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"cmt", "cm", "c"},
	Short:   "Commit container state.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDefFromCommand(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.Commit())
				return
			}
		}
		handleError(fmt.Errorf("can only commit applications"))
	},
}

var containerAppDeleteCommitCmd = &cobra.Command{
	Use:     "delete_commit",
	Aliases: []string{"dc", "dcmt", "dcommit", "dcm"},
	Short:   "Delete container comit.",
	Run: func(cmd *cobra.Command, args []string) {
		proj, err := getProject(true)
		handleError(err)
		d, err := getDefFromCommand(containerCmd, proj)
		handleError(err)
		switch d.(type) {
		case def.App:
			{
				c := proj.NewContainer(d)
				handleError(c.DeleteCommit())
				return
			}
		}
		handleError(fmt.Errorf("can only commit applications"))
	},
}

var containerLogsCmd = &cobra.Command{
	Use:   "logs [-f follow]",
	Short: "Display logs for container.",
	Run: func(cmd *cobra.Command, args []string) {
		output.Enable = false
		proj, err := getProject(true)
		handleError(err)
		followFlag := cmd.Flags().Lookup("follow")
		hasFollow := followFlag != nil && followFlag.Value.String() != "false"
		d, err := getDefFromCommand(containerCmd, proj)
		handleError(err)
		handleError(proj.NewContainer(d).LogStdout(hasFollow))
		if hasFollow {
			select {}
		}
	},
}

var containerCopyCmd = &cobra.Command{
	Use:     "copy source destination",
	Aliases: []string{"cp"},
	Short:   "Copy file to and from container.",
	Run: func(cmd *cobra.Command, args []string) {
		// ensure two args provided
		if len(args) != 2 {
			handleError(fmt.Errorf("unexpected number of arguements"))
		}
		// fetch project
		proj, err := getProject(true)
		handleError(err)

		// log output
		done := output.Duration(fmt.Sprintf("Copy %s to %s.", args[0], args[1]))

		// determine path local vs container
		parsePath := func(path string) (string, string) {
			pSplit := strings.Split(path, ":")
			if len(pSplit) == 1 {
				// local file
				return "", path
			}
			// container apth
			return pSplit[0], strings.Join(pSplit[1:], ":")
		}

		// source
		srcService, srcPath := parsePath(args[0])
		var srcReader *os.File
		if srcService == "" {
			// local
			var err error
			srcReader, err = os.Open(srcPath)
			handleError(err)
			defer srcReader.Close()
		} else {
			// container
			def, err := getDef(srcService, proj)
			handleError(err)
			cont := proj.NewContainer(def)
			srcReader, err = ioutil.TempFile(os.TempDir(), "pcc-cp-")
			handleError(err)
			defer func() {
				name := srcReader.Name()
				srcReader.Close()
				os.Remove(name)
			}()
			handleError(cont.Download(srcPath, srcReader))
			srcReader.Seek(0, 0)
		}
		if srcReader == nil {
			handleError(fmt.Errorf("invalid source reader"))
		}

		// dest
		destService, destPath := parsePath(args[1])
		if destService == "" {

			// if dest path is a directory then use source filename as dest filename
			stat, err := os.Stat(destPath)
			handleError(err)
			if stat.IsDir() {
				destPath = filepath.Join(destPath, filepath.Base(srcPath))
			}

			// local
			dstWriter, err := os.Create(destPath)
			handleError(err)
			defer dstWriter.Close()
			_, err = io.Copy(dstWriter, srcReader)
			handleError(err)

		} else {
			// container
			def, err := getDef(destService, proj)
			handleError(err)
			cont := proj.NewContainer(def)
			handleError(cont.Upload(destPath, srcReader))
		}
		done()
	},
}

func init() {
	containerShellCmd.PersistentFlags().Bool("root", false, "shell as root")
	containerLogsCmd.Flags().BoolP("follow", "f", false, "follow logs")
	containerCmd.PersistentFlags().StringP("service", "s", "", "name of service/application/worker")
	containerCmd.AddCommand(containerAppDeployCmd)
	containerCmd.AddCommand(containerAppPostDeployCmd)
	containerCmd.AddCommand(containerShellCmd)
	containerCmd.AddCommand(containerAppCommitCmd)
	containerCmd.AddCommand(containerAppDeleteCommitCmd)
	containerCmd.AddCommand(containerLogsCmd)
	containerCmd.AddCommand(containerCopyCmd)
	RootCmd.AddCommand(containerCmd)
}
