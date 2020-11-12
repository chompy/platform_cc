package api

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"golang.org/x/crypto/ssh/terminal"
)

// ShellContainer - access shell inside given container
func (d *dockerClient) ShellContainer(id string, user string, command []string) error {
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
		return err
	}
	defer hresp.Close()
	// don't create interactive shell for stdin
	if hasStdin {
		_, err = io.Copy(hresp.Conn, os.Stdin)
		return err
	}
	// create interactive shell
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()
	go func() { io.Copy(hresp.Conn, os.Stdin) }()
	io.Copy(os.Stdout, hresp.Reader)
	return nil
}