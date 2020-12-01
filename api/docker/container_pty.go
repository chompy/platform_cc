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

package docker

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/ztrue/tracerr"
	"golang.org/x/crypto/ssh/terminal"
)

// resizeShell resizes the given Docker process to match the current terminal.
func (d *Client) resizeShell(execID string) error {
	w, h, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return tracerr.Wrap(err)
	}
	resizeOpts := types.ResizeOptions{
		Width:  uint(w),
		Height: uint(h),
	}
	err = d.cli.ContainerExecResize(
		context.Background(),
		execID,
		resizeOpts,
	)
	return tracerr.Wrap(err)
}

// ShellContainer creates an interactive shell in given container.
func (d *Client) ShellContainer(id string, user string, command []string) error {
	// check stdin
	fi, _ := os.Stdin.Stat()
	hasStdin := false
	if fi.Size() > 0 {
		hasStdin = true
	}
	// open shell
	resp, err := d.cli.ContainerExecCreate(
		context.Background(),
		id,
		types.ExecConfig{
			User:         user,
			Tty:          !hasStdin,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          command,
		},
	)
	hresp, err := d.cli.ContainerExecAttach(
		context.Background(),
		resp.ID,
		types.ExecConfig{
			User:         user,
			Tty:          !hasStdin,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          command,
		},
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer hresp.Close()
	// don't create interactive shell if stdin already exists
	if hasStdin {
		_, err = io.Copy(hresp.Conn, os.Stdin)
		return tracerr.Wrap(err)
	}
	// create interactive shell
	// handle resizing
	err = d.resizeShell(resp.ID)
	if err != nil {
		return tracerr.Wrap(err)
	}
	go func() {
		for {
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGWINCH)
			s := <-sigc
			switch s {
			case syscall.SIGWINCH:
				{
					d.resizeShell(resp.ID)
					break
				}
			}
		}
	}()
	// make raw
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()
	// read/write connection to stdin and stdout
	go func() { io.Copy(hresp.Conn, os.Stdin) }()
	io.Copy(os.Stdout, hresp.Reader)
	return nil
}
