package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/ztrue/tracerr"
)

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
	// don't create interactive shell for stdin
	if hasStdin {
		_, err = io.Copy(hresp.Conn, os.Stdin)
		return tracerr.Wrap(err)
	}
	// use docker cli to shell for best experience
	// TODO figure out how to improve direct terminal access
	cmd := exec.Command(
		"sh", "-c",
		fmt.Sprintf("docker exec --user %s -i -t %s sh -c '%s'", user, id, strings.Join(command, " ")),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil
	}
	// create interactive shell
	/*oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()
	go func() { io.Copy(hresp.Conn, os.Stdin) }()
	io.Copy(os.Stdout, hresp.Reader)*/
	return nil
}
