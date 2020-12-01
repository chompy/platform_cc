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
	"time"

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

// handleResizeShell resizes the given Docker process anytime the current terminal is resized.
func (d *Client) handleResizeShell(execID string) error {
	cw, ch, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return tracerr.Wrap(err)
	}
	for range time.Tick(time.Millisecond * 100) {
		w, h, err := terminal.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			return tracerr.Wrap(err)
		}
		if cw != w || ch != h {
			if err := d.resizeShell(execID); err != nil {
				return tracerr.Wrap(err)
			}
			cw = w
			ch = h
		}
	}
	return nil
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
	if err := d.resizeShell(resp.ID); err != nil {
		return tracerr.Wrap(err)
	}
	/*if err != nil {
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
	}()*/
	go d.handleResizeShell(resp.ID)

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
