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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
	"golang.org/x/term"
)

type exitResult struct {
	code  int
	error error
}

// resizeShell resizes the given Docker process to match the current term.
func (d Docker) resizeShell(execID string) error {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return errors.WithStack(err)
	}
	resizeOpts := types.ResizeOptions{
		Width:  uint(w),
		Height: uint(h),
	}
	err = d.client.ContainerExecResize(
		context.Background(),
		execID,
		resizeOpts,
	)
	return errors.WithStack(err)
}

// handleResizeShell resizes the given Docker process anytime the current terminal is resized.
func (d Docker) handleResizeShell(execID string) error {
	cw, ch, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return errors.WithStack(err)
	}
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	done := make(chan bool)
	for {
		select {
		case <-done:
			{
				return nil
			}
		case <-ticker.C:
			{
				w, h, err := term.GetSize(int(os.Stdin.Fd()))
				if err != nil {
					return errors.WithStack(err)
				}
				if cw != w || ch != h {
					if err := d.resizeShell(execID); err != nil {
						return errors.WithStack(err)
					}
					cw = w
					ch = h
				}
			}
		}
	}
}

// ContainerShell creates an interactive shell in given container.
func (d Docker) ContainerShell(id string, user string, cmd []string, stdin io.Reader) (int, error) {
	// ensure container is running
	if status, _ := d.ContainerStatus(id); !status.Running {
		return -1, errors.Wrapf(ErrContainerNotRunning, "container %s is not running", id)
	}
	// check stdin
	hasStdin := true
	if stdin == nil {
		fi, _ := os.Stdin.Stat()
		hasStdin = fi.Mode()&os.ModeDevice == 0
		stdin = os.Stdin
	}
	// open shell
	execConfig := types.ExecConfig{
		User:         user,
		Tty:          !hasStdin,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	output.LogDebug(
		fmt.Sprintf("Docker open shell. (Container ID %s)", id), execConfig,
	)
	resp, err := d.client.ContainerExecCreate(
		context.Background(),
		id,
		execConfig,
	)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	execConfig = types.ExecConfig{
		User:         user,
		Tty:          !hasStdin,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	output.LogDebug(
		fmt.Sprintf("Docker container attach. (Exec ID %s)", resp.ID), execConfig,
	)
	hresp, err := d.client.ContainerExecAttach(
		context.Background(),
		resp.ID,
		execConfig,
	)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	defer hresp.Close()
	// don't create interactive shell if stdin already exists
	if hasStdin {
		output.LogDebug("Disable interactive shell, stdin present.", nil)
		// process stdin
		exit := make(chan exitResult)
		go func(exit chan exitResult) {
			var n int64
			var err error
			for {
				n, err = io.Copy(hresp.Conn, stdin)
				if err == nil {
					code, err := d.checkCommandExec(resp.ID)
					exit <- exitResult{code, errors.WithStack(err)}
					return
				}
				if strings.Contains(err.Error(), "broken pipe") {
					output.LogDebug(fmt.Sprintf("Copy stdin broken pipe after writing %d bytes.", n), resp.ID)
					hresp, err = d.client.ContainerExecAttach(
						context.Background(),
						resp.ID,
						execConfig,
					)
					if err != nil {
						exit <- exitResult{-1, errors.WithStack(err)}
						return
					}
					continue
				}
				code, err := d.checkCommandExec(resp.ID)
				exit <- exitResult{code, errors.WithStack(err)}
				return
			}
		}(exit)
		// pipe stdout
		scanner := bufio.NewScanner(hresp.Reader)
		for scanner.Scan() {
			t := scanner.Bytes()
			output.WriteStdout(string(t) + "\n")
		}
		res := <-exit
		return res.code, errors.WithStack(res.error)
	}
	// create interactive shell
	// handle resizing
	d.resizeShell(resp.ID)
	go d.handleResizeShell(resp.ID)
	// make raw
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return -1, errors.WithStack(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
	// read/write connection to stdin and stdout
	go func() { io.Copy(hresp.Conn, stdin) }()
	io.Copy(os.Stdout, hresp.Reader)
	code, err := d.checkCommandExec(resp.ID)
	return code, errors.WithStack(err)
}
