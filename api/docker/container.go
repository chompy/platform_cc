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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// StartContainer creates and starts a Docker container.
func (d *Client) StartContainer(c ContainerConfig) error {
	done := output.Duration(
		fmt.Sprintf(
			"Start Docker container for %s '%s.'",
			c.ObjectType.TypeName(),
			c.ObjectName,
		),
	)
	// get mounts
	mounts := make([]mount.Mount, 0)
	for k, v := range c.Volumes {
		if k == "" || v == "" {
			continue
		}
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: k,
			Target: v,
		})
	}
	for k, v := range c.Binds {
		if k == "" || v == "" {
			continue
		}
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: k,
			Target: v,
		})
	}
	// get port mappings
	exposedPorts, portBinding, err := nat.ParsePortSpecs(c.Ports)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// create container
	resp, err := d.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        c.Image,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          c.GetCommand(),
			Env:          c.GetEnv(),
			WorkingDir:   c.WorkingDir,
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			AutoRemove:   true,
			Privileged:   true,
			Mounts:       mounts,
			PortBindings: portBinding,
		},
		&network.NetworkingConfig{},
		c.GetContainerName(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			return nil
		}
		return tracerr.Wrap(err)
	}
	for _, w := range resp.Warnings {
		output.Warn(w)
	}
	// attach container to project network
	err = d.cli.NetworkConnect(
		context.Background(),
		globalNetworkName,
		resp.ID,
		nil,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// start container
	err = d.cli.ContainerStart(
		context.Background(),
		resp.ID,
		types.ContainerStartOptions{},
	)
	if err == nil {
		done()
	}
	return tracerr.Wrap(err)
}

// GetProjectContainers gets a list of active containers for given project.
func (d *Client) GetProjectContainers(pid string) ([]types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(dockerNamingPrefix+"*", pid))
	return d.cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: filterArgs,
		},
	)
}

// GetAllContainers gets a list of all active PCC containers.
func (d *Client) GetAllContainers() ([]types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	return d.cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: filterArgs,
		},
	)
}

// DeleteProjectContainers deletes all containers related to a project.
func (d *Client) DeleteProjectContainers(pid string) error {
	containers, err := d.GetProjectContainers(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteContainers(containers)
}

// DeleteAllContainers deletes all PCC containers.
func (d *Client) DeleteAllContainers() error {
	containers, err := d.GetAllContainers()
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteContainers(containers)
}

// deleteContainers deletes all provided containers.
func (d *Client) deleteContainers(containers []types.Container) error {
	timeout := 30 * time.Second
	ch := make(chan error)
	msgs := make([]string, 0)
	for _, c := range containers {
		//log.Printf("Delete Docker container '%s.'", c.Names[0])
		msgs = append(msgs, strings.Trim(c.Names[0], "/"))
		go func(cid string) {
			if err := d.cli.ContainerStop(
				context.Background(),
				cid,
				&timeout,
			); err != nil {
				ch <- err
			}
			ch <- nil
		}(c.ID)
	}
	// output progress
	done := output.Duration("Delete Docker containers.")
	prog := output.Progress(msgs)

	// wait for deletion
	for i := range containers {
		err := <-ch
		if err != nil {
			prog(i, output.ProgressMessageError)
			return tracerr.Wrap(err)
		}
		prog(i, output.ProgressMessageDone)
	}
	done()
	return nil
}

// RunContainerCommand runs a command in a container.
func (d *Client) RunContainerCommand(id string, user string, cmd []string, out io.Writer) error {
	execConfig := types.ExecConfig{
		User:         user,
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	resp, err := d.cli.ContainerExecCreate(
		context.Background(),
		id,
		execConfig,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	hresp, err := d.cli.ContainerExecAttach(
		context.Background(),
		resp.ID,
		execConfig,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if out != nil {
		_, err = io.Copy(out, hresp.Reader)
	}
	return tracerr.Wrap(err)
}

// UploadDataToContainer uploads data to container as a file.
func (d *Client) UploadDataToContainer(id string, path string, r io.Reader) error {
	// TODO this is not the best way to upload a file to the container but it's the only
	// one that seems to work for now
	// read data as bytes
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// gzip data stream
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return tracerr.Wrap(err)
	}
	if err := zw.Close(); err != nil {
		return tracerr.Wrap(err)
	}
	// convert to base64 string
	dataB64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	if err := d.RunContainerCommand(
		id,
		"root",
		[]string{"sh", "-c", "echo '" + dataB64 + "' | base64 -d | gunzip -c > " + path},
		nil,
	); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// GetContainerIP gets the IP address of the given container.
func (d *Client) GetContainerIP(id string) (string, error) {
	data, err := d.cli.ContainerInspect(
		context.Background(),
		id,
	)
	if err != nil {
		return "", err
	}
	for name, network := range data.NetworkSettings.Networks {
		if name == globalNetworkName {
			return network.IPAddress, nil
		}
	}
	return "", tracerr.Wrap(fmt.Errorf("network not found for container '%s'", id))

}

// PullImage pulls the latest image for the given container.
func (d *Client) PullImage(c ContainerConfig) error {
	done := output.Duration(
		fmt.Sprintf("Pull Docker container image for '%s.'", c.GetContainerName()),
	)
	r, err := d.cli.ImagePull(
		context.Background(),
		c.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		return tracerr.Wrap(err)
	}

	defer r.Close()
	b := make([]byte, 1024)
	for {
		n, _ := r.Read(b)
		if n == 0 {
			break
		}
	}
	done()
	return nil
}

// PullImages pulls images for multiple containers.
func (d *Client) PullImages(containerConfigs []ContainerConfig) error {

	// get list of images and prepare progress output
	images := make([]string, 0)
	msgs := make([]string, 0)
	for _, c := range containerConfigs {
		hasImage := false
		for _, image := range images {
			if image == c.Image {
				hasImage = true
				break
			}
		}
		if !hasImage {
			images = append(images, c.Image)
			msgs = append(msgs, c.GetHumanName())
		}
	}
	prog := output.Progress(msgs)
	// simultaneously pull images
	outputEnabled := output.Enable
	var wg sync.WaitGroup
	for i, c := range containerConfigs {
		wg.Add(1)
		go func(i int, c ContainerConfig) {
			defer wg.Done()
			defer func() { output.Enable = outputEnabled }()
			output.Enable = false
			if err := d.PullImage(c); err != nil {
				prog(i, output.ProgressMessageError)
				output.Warn(err.Error())
			}
			prog(i, output.ProgressMessageDone)
		}(i, c)

	}
	wg.Wait()
	return nil
}
