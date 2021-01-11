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
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
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
func (d MainClient) StartContainer(c ContainerConfig) error {
	if d.isContainerRunning(c.GetContainerName()) {
		return nil
	}
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
	cConfig := &container.Config{
		Image:        c.Image,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          c.GetCommand(),
		Env:          c.GetEnv(),
		WorkingDir:   c.WorkingDir,
		ExposedPorts: exposedPorts,
	}
	cHostConfig := &container.HostConfig{
		AutoRemove:   true,
		Privileged:   true,
		Mounts:       mounts,
		PortBindings: portBinding,
	}
	output.LogDebug(fmt.Sprintf("Container create. (Name %s)", c.GetContainerName()), []interface{}{cConfig, cHostConfig})
	resp, err := d.cli.ContainerCreate(
		context.Background(),
		cConfig,
		cHostConfig,
		&network.NetworkingConfig{},
		c.GetContainerName(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			return nil
		}
		return tracerr.Wrap(err)
	}
	output.LogDebug("Container created.", resp)
	// attach container to project network
	if err := d.cli.NetworkConnect(
		context.Background(),
		globalNetworkName,
		resp.ID,
		nil,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// start container
	if err := d.cli.ContainerStart(
		context.Background(),
		resp.ID,
		types.ContainerStartOptions{},
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// GetProjectContainers gets a list of active containers for given project.
func (d MainClient) GetProjectContainers(pid string) ([]types.Container, error) {
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
func (d MainClient) GetAllContainers() ([]types.Container, error) {
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
func (d MainClient) DeleteProjectContainers(pid string) error {
	containers, err := d.GetProjectContainers(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteContainers(containers)
}

// DeleteAllContainers deletes all PCC containers.
func (d MainClient) DeleteAllContainers() error {
	containers, err := d.GetAllContainers()
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteContainers(containers)
}

// deleteContainers deletes all provided containers.
func (d MainClient) deleteContainers(containers []types.Container) error {
	timeout := 30 * time.Second
	output.LogDebug("Delete containers.", containers)
	// output progress
	msgs := make([]string, len(containers))
	for i := range containers {
		msgs[i] = strings.Trim(containers[i].Names[0], "/")
	}
	done := output.Duration("Delete Docker containers.")
	prog := output.Progress(msgs)
	// itterate containers and stop
	var wg sync.WaitGroup
	for i := range containers {
		wg.Add(1)
		go func(cid string, i int) {
			defer wg.Done()
			if err := d.cli.ContainerStop(
				context.Background(),
				cid,
				&timeout,
			); err != nil {
				prog(i, output.ProgressMessageError)
				output.LogError(err)
			}
			prog(i, output.ProgressMessageDone)
		}(containers[i].ID, i)
	}
	wg.Wait()
	done()
	return nil
}

// RunContainerCommand runs a command in a container.
func (d MainClient) RunContainerCommand(id string, user string, cmd []string, out io.Writer) error {
	execConfig := types.ExecConfig{
		User:         user,
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	output.LogDebug(fmt.Sprintf("Container exec create. (Container ID %s)", id), execConfig)
	resp, err := d.cli.ContainerExecCreate(
		context.Background(),
		id,
		execConfig,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	output.LogDebug("Container exec created.", resp.ID)
	hresp, err := d.cli.ContainerExecAttach(
		context.Background(),
		resp.ID,
		execConfig,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// get command stdout
	var buf bytes.Buffer
	var mWriter io.Writer
	mWriter = &buf
	if out != nil {
		mWriter = io.MultiWriter(&buf, out)
	}
	if _, err := io.Copy(mWriter, hresp.Reader); err != nil {
		return tracerr.Wrap(err)
	}
	output.LogDebug("Container exec finished.", string(buf.Bytes()))
	return nil
}

// UploadDataToContainer uploads data to container as a file.
func (d MainClient) UploadDataToContainer(id string, path string, r io.Reader) error {
	// log debug
	output.LogDebug(
		"Upload to container.",
		map[string]interface{}{
			"container_id": id,
			"path":         path,
		},
	)
	// gzip data stream
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := io.Copy(zw, r); err != nil {
		return tracerr.Wrap(err)
	}
	if err := zw.Close(); err != nil {
		return tracerr.Wrap(err)
	}
	// push file to container via stdin
	return d.ShellContainer(
		id, "root",
		[]string{"sh", "-c", fmt.Sprintf("cat | gunzip > %s", path)},
		&buf,
	)
}

// GetContainerIP gets the IP address of the given container.
func (d MainClient) GetContainerIP(id string) (string, error) {
	data, err := d.cli.ContainerInspect(
		context.Background(),
		id,
	)
	if err != nil {
		return "", err
	}
	output.LogDebug(
		fmt.Sprintf("Inspected container (%s) for IP address.", id),
		data.NetworkSettings,
	)
	for name, network := range data.NetworkSettings.Networks {
		if name == globalNetworkName {
			return network.IPAddress, nil
		}
	}
	return "", tracerr.Wrap(fmt.Errorf("network not found for container '%s'", id))

}

// isContainerRunning returns true if given container is running.
func (d MainClient) isContainerRunning(id string) bool {
	data, err := d.cli.ContainerInspect(
		context.Background(),
		id,
	)
	if err != nil {
		return false
	}
	if data.State == nil {
		return false
	}
	return data.State.Running
}

// PullImage pulls the latest image for the given container.
func (d MainClient) PullImage(c ContainerConfig) error {
	done := output.Duration(
		fmt.Sprintf("Pull Docker container image for '%s.'", c.GetContainerName()),
	)
	output.LogDebug(
		"Pull container image.",
		map[string]interface{}{
			"container_id": c.GetContainerName(),
			"image":        c.Image,
		},
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
func (d MainClient) PullImages(containerConfigs []ContainerConfig) error {
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
				output.LogError(err)
			}
			prog(i, output.ProgressMessageDone)
		}(i, c)
	}
	wg.Wait()
	return nil
}

// ContainerLog logs container output.
func (d MainClient) ContainerLog(id string) {
	output.LogInfo(fmt.Sprintf("Read logs for container '%s.'", id))
	go func() {
		res, err := d.cli.ContainerLogs(
			context.Background(),
			id,
			types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     true,
			},
		)
		if err != nil {
			output.LogError(err)
			return
		}
		scanner := bufio.NewScanner(res)
		defer res.Close()
		for {
			for scanner.Scan() {
				output.LogDebug(fmt.Sprintf("[%s] %s", id, scanner.Text()), nil)
			}
			if err := scanner.Err(); err != nil {
				output.LogError(err)
			}
		}
	}()
}
