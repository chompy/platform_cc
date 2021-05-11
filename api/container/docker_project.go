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
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/mount"
	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// ProjectStop stops all running Docker containers for given project.
func (d Docker) ProjectStop(pid string) error {
	containers, err := d.listProjectContainers(pid)
	if err != nil {
		return errors.WithStack(err)
	}
	return d.deleteContainers(containers)
}

// ProjectPurge deletes all Docker resources for given project.
func (d Docker) ProjectPurge(pid string) error {
	// stop
	if err := d.ProjectStop(pid); err != nil {
		return errors.WithStack(err)
	}
	time.Sleep(time.Second)
	// delete volumes
	vols, err := d.listProjectVolumes(pid)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.deleteVolumes(vols); err != nil {
		return errors.WithStack(err)
	}
	time.Sleep(time.Second)
	// delete images
	images, err := d.listProjectImages(pid)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.deleteImages(images); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ProjectPurgeSlot deletes all Docker resources for given project slot.
func (d Docker) ProjectPurgeSlot(pid string, slot int) error {
	// stop
	if err := d.ProjectStop(pid); err != nil {
		return errors.WithStack(err)
	}
	// delete volumes
	vols, err := d.listProjectSlotVolumes(pid, slot)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.deleteVolumes(vols); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ProjectCopySlot copies volumes in given slot to another slot.
func (d Docker) ProjectCopySlot(pid string, sourceSlot int, destSlot int) error {
	// can't have same slots
	if sourceSlot == destSlot {
		return errors.Wrap(ErrInvalidSlot, "source and destination slots cannot be the same")
	}
	// purge dest slot in prep for copy
	d.ProjectPurgeSlot(pid, destSlot)
	// log copy start
	done := output.Duration(fmt.Sprintf("Copy slot %d to %d.", sourceSlot, destSlot))
	// get list of source slot volumes
	volList, err := d.listProjectSlotVolumes(pid, sourceSlot)
	if err != nil {
		return errors.WithStack(err)
	}
	// create mounts for container, source and dest volumes for copy
	mounts := make([]mount.Mount, 0)
	for _, vol := range volList.Volumes {
		mounts = append(
			mounts,
			mount.Mount{
				Type:   mount.TypeVolume,
				Source: vol.Name,
				Target: "/mnt/src/" + volumeStripSlot(vol.Name),
			},
		)
		mounts = append(
			mounts,
			mount.Mount{
				Type:   mount.TypeVolume,
				Source: volumeWithSlot(vol.Name, destSlot),
				Target: "/mnt/dest/" + volumeStripSlot(vol.Name),
			},
		)
	}
	// create container
	cConfig := &container.Config{
		Image:        "busybox",
		Cmd:          []string{"sh", "-c", "cp -rv /mnt/src/* /mnt/dest/"},
		AttachStdout: true,
	}
	cHostConfig := &container.HostConfig{
		AutoRemove: true,
		Mounts:     mounts,
	}
	output.LogDebug("Slot copy container create.", []interface{}{cConfig, cHostConfig})
	resp, err := d.client.ContainerCreate(
		context.Background(),
		cConfig,
		cHostConfig,
		&network.NetworkingConfig{},
		fmt.Sprintf(containerNamingPrefix, pid)+"slotcopy",
	)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			return nil
		}
		return errors.WithStack(convertDockerError(err))
	}
	output.LogDebug("Slot copy container created.", resp)
	// start container
	if err := d.client.ContainerStart(
		context.Background(),
		resp.ID,
		types.ContainerStartOptions{},
	); err != nil {
		return errors.WithStack(convertDockerError(err))
	}
	cleanup := func() {
		timeout := time.Second
		d.client.ContainerStop(
			context.Background(), resp.ID, &timeout,
		)
	}
	defer cleanup()
	// attach container
	attachResp, err := d.client.ContainerAttach(
		context.Background(),
		resp.ID,
		types.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		return errors.WithStack(convertDockerError(err))
	}
	defer attachResp.Close()
	// capture interupt to stop container
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		<-c
		output.Info("Interupt detected, stopping.")
		cleanup()
		attachResp.Close()
		os.Exit(1)
	}()
	// pipe to stdout
	if _, err := io.Copy(os.Stdout, attachResp.Reader); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}
