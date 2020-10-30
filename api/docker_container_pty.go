package api

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"golang.org/x/crypto/ssh/terminal"
)

// ShellContainer - access shell inside given container
func (d *dockerClient) ShellContainer(id string) error {

	resp, err := d.cli.ContainerExecCreate(
		context.Background(),
		id,
		types.ExecConfig{
			User:         "web",
			Tty:          true,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{"bash"},
		},
	)

	hresp, err := d.cli.ContainerExecAttach(
		context.Background(),
		resp.ID,
		types.ExecConfig{
			User:         "web",
			Tty:          true,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{"bash"},
		},
	)
	if err != nil {
		return err
	}
	defer hresp.Close()
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
	go func() { io.Copy(hresp.Conn, os.Stdin) }()
	io.Copy(os.Stdout, hresp.Reader)
	return nil
}
