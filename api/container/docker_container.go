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
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"path"
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

const dockerCommitTagPrefix = "pcc.local/build:"
const containerStopTimeout = 10

// ContainerStart starts a Docker container.
func (d Docker) ContainerStart(c Config) error {
	// ensure not already running
	if status, _ := d.ContainerStatus(c.GetContainerName()); status.Running {
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
	// volumes
	for k, v := range c.Volumes {
		if k == "" || v == "" {
			continue
		}
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: volumeWithSlot(getMountName(c.ProjectID, k, c.ObjectType), c.Slot),
			Target: v,
		})
	}
	// binds
	for k, v := range c.Binds {
		if k == "" || v == "" {
			continue
		}
		// use nfs if macos
		if isMacOS() && c.EnableOSXNFS {
			if err := d.createNFSVolume(
				c.ProjectID, path.Base(k), k, c.ObjectType,
			); err != nil {
				return tracerr.Wrap(err)
			}
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeVolume,
				Source: getMountName(c.ProjectID, path.Base(k), c.ObjectType),
				Target: v,
			})
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
	// create network
	if err := d.createNetwork(); err != nil {
		return tracerr.Wrap(err)
	}
	// check for committed image
	image := fmt.Sprintf("%s%s", dockerCommitTagPrefix, c.GetContainerName())
	if c.NoCommit || !d.hasImage(image) {
		image = c.Image
	}
	// create container
	cConfig := &container.Config{
		Image:        image,
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
	resp, err := d.client.ContainerCreate(
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
	if err := d.client.NetworkConnect(
		context.Background(),
		dockerNetworkName,
		resp.ID,
		nil,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// start container
	if err := d.client.ContainerStart(
		context.Background(),
		resp.ID,
		types.ContainerStartOptions{},
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// ContainerCommand runs a command inside a Docker container.
func (d Docker) ContainerCommand(id string, user string, cmd []string, out io.Writer) error {
	execConfig := types.ExecConfig{
		User:         user,
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	output.LogDebug(fmt.Sprintf("Container exec create. (Container ID %s)", id), execConfig)
	resp, err := d.client.ContainerExecCreate(
		context.Background(),
		id,
		execConfig,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	output.LogDebug("Container exec created.", resp.ID)
	hresp, err := d.client.ContainerExecAttach(
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

// ContainerStatus returns status of Docker container.
func (d Docker) ContainerStatus(id string) (Status, error) {
	data, err := d.client.ContainerInspect(
		context.Background(),
		id,
	)
	if err != nil {
		return Status{Running: false, Slot: 0}, tracerr.Wrap(err)
	}
	ipAddress := ""
	for name, network := range data.NetworkSettings.Networks {
		if name == dockerNetworkName {
			ipAddress = network.IPAddress
			break
		}
	}
	slot := 1
	for _, m := range data.Mounts {
		if m.Type == mount.TypeVolume {
			slot = volumeGetSlot(m.Name)
			break
		}
	}
	config := containerConfigFromName(data.Name)
	return Status{
		ID:         data.ID,
		Name:       config.ObjectName,
		ObjectType: config.ObjectType,
		Image:      data.Image,
		Type:       d.serviceTypeFromImage(data.Image),
		Committed:  d.hasCommit(data.Name),
		ProjectID:  config.ProjectID,
		Running:    data.State.Running,
		IPAddress:  ipAddress,
		Slot:       slot,
	}, nil
}

// ContainerUpload uploads a file to a Docker container.
func (d Docker) ContainerUpload(id string, path string, r io.Reader) error {
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
	return d.ContainerShell(
		id, "root",
		[]string{"sh", "-c", fmt.Sprintf("cat | gunzip > %s", path)},
		&buf,
	)
}

// ContainerLog returns a reader containing log data for a Docker container.
func (d Docker) ContainerLog(id string, follow bool) (io.ReadCloser, error) {
	output.LogDebug(fmt.Sprintf("Read logs for container '%s.'", id), nil)
	return d.client.ContainerLogs(
		context.Background(),
		id,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     follow,
		},
	)
}

// ContainerCommit stores a Docker container's state as an image.
func (d Docker) ContainerCommit(id string) error {
	// check container state
	data, err := d.client.ContainerInspect(
		context.Background(),
		id,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if !data.State.Running {
		return tracerr.Errorf("container %s is not running", id)
	}
	done := output.Duration(
		fmt.Sprintf(
			"Commit container '%s.'",
			id,
		),
	)
	// commit image
	idResp, err := d.client.ContainerCommit(
		context.Background(),
		data.ID,
		types.ContainerCommitOptions{},
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// tag image
	if err := d.client.ImageTag(
		context.Background(),
		idResp.ID,
		fmt.Sprintf("%s%s", dockerCommitTagPrefix, strings.Trim(data.Name, "/")),
	); err != nil {
		d.client.ImageRemove(
			context.Background(), idResp.ID, types.ImageRemoveOptions{Force: true},
		)
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// ContainerDeleteCommit deletes Docker image for given container.
func (d Docker) ContainerDeleteCommit(id string) error {
	done := output.Duration(
		fmt.Sprintf(
			"Delete committed container image '%s.'",
			id,
		),
	)
	// ensure container isn't running
	s, _ := d.ContainerStatus(id)
	if s.Running {
		return tracerr.New("cannot delete commit for running container")
	}
	// check for committed image
	image := fmt.Sprintf("%s%s", dockerCommitTagPrefix, id)
	if !d.hasImage(image) {
		return tracerr.New(fmt.Sprintf("container %s has no comitted image", id))
	}
	// delete
	if _, err := d.client.ImageRemove(
		context.Background(),
		image,
		types.ImageRemoveOptions{Force: true},
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// listProjectContainers gets a list of active containers for given project.
func (d Docker) listProjectContainers(pid string) ([]types.Container, error) {
	output.LogDebug(fmt.Sprintf("List containers for project '%s.'", pid), nil)
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(containerNamingPrefix+"*", pid))
	return d.client.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: filterArgs,
		},
	)
}

// listAllContainers gets a list of all active Platform.CC containers.
func (d Docker) listAllContainers() ([]types.Container, error) {
	output.LogDebug("List all Platform.CC containers.", nil)
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	return d.client.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: filterArgs,
		},
	)
}

// deleteContainers deletes all provided containers.
func (d Docker) deleteContainers(containers []types.Container) error {
	//timeout := 30 * time.Second
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
			// issue shutdown command to container, containers are set to be auto deleted
			// once they exit
			c := make(chan error, 1)
			go func() {
				c <- d.ContainerCommand(
					cid, "root", []string{"sh", "-c", "/etc/platform/shutdown || true && rm -f /routes/* && nginx -s stop || true"}, nil,
				)
			}()
			select {
			case err := <-c:
				{
					if err != nil {
						prog(i, output.ProgressMessageError, nil, nil)
						output.LogError(err)
					}
					prog(i, output.ProgressMessageDone, nil, nil)
					return
				}
			case <-time.After(time.Second * containerStopTimeout):
				{
					output.LogDebug(
						fmt.Sprintf("Delete container '%s' timed out after %d seconds, forcing delete.", cid, containerStopTimeout),
						nil,
					)
					timeout := time.Second * containerStopTimeout
					if err := d.client.ContainerStop(
						context.Background(),
						cid,
						&timeout,
					); err != nil {
						prog(i, output.ProgressMessageError, nil, nil)
						output.LogError(err)
						return
					}
					prog(i, output.ProgressMessageDone, nil, nil)
					return
				}
			}
		}(containers[i].ID, i)
	}
	wg.Wait()
	done()
	return nil
}

// hasCommit return true if a committed version of the given container exists.
func (d Docker) hasCommit(name string) bool {
	name = strings.TrimPrefix(name, "/")
	image := fmt.Sprintf("%s%s", dockerCommitTagPrefix, name)
	return d.hasImage(image)
}
