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

package container

import (
	"io"
)

// Interface defines methods used to interact with container.
type Interface interface {
	ContainerStart(c Config) error
	ContainerCommand(id string, user string, cmd []string, out io.Writer) error
	ContainerShell(id string, user string, cmd []string, stdin io.Reader) error
	ContainerStatus(id string) (Status, error)
	ContainerUpload(id string, path string, r io.Reader) error
	ContainerDownload(id string, path string, w io.Writer) error
	ContainerLog(id string, follow bool) (io.ReadCloser, error)
	ContainerCommit(id string) error
	ContainerDeleteCommit(id string) error
	ImagePull(c []Config) error
	ProjectStop(pid string) error
	ProjectPurge(pid string) error
	ProjectPurgeSlot(pid string, slot int) error
	ProjectCopySlot(pid string, sourceSlot int, destSlot int) error
	AllStop() error
	AllPurge(deleteGlobalVolumes bool) error
	AllStatus() ([]Status, error)
}
