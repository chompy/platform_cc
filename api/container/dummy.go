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
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
)

// DummyTracker trackes
type DummyTracker struct {
	Volumes    []string
	Containers []string
	Sync       sync.Mutex
}

// Dummy defines the dummy container handler.
type Dummy struct {
	Tracker *DummyTracker
}

// NewDummy creates a new dummy handler.
func NewDummy() Dummy {
	return Dummy{
		Tracker: &DummyTracker{
			Volumes:    make([]string, 0),
			Containers: make([]string, 0),
		},
	}
}

func (d Dummy) hasContainer(id string) bool {
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	for _, c := range d.Tracker.Containers {
		if c == id {
			return true
		}
	}
	return false
}

func (d Dummy) hasVolume(id string) bool {
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	for _, v := range d.Tracker.Volumes {
		if v == id {
			return true
		}
	}
	return false
}

// ContainerStart starts dummy container.
func (d Dummy) ContainerStart(c Config) error {
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	d.Tracker.Containers = append(d.Tracker.Containers, c.GetContainerName())
	for k := range c.Volumes {
		d.Tracker.Volumes = append(d.Tracker.Volumes, c.ProjectID+"_"+k)
	}
	return nil
}

// ContainerCommand runs dummy command.
func (d Dummy) ContainerCommand(id string, user string, cmd []string, out io.Writer) error {
	if !d.hasContainer(id) {
		return fmt.Errorf("container %s not running", id)
	}
	return nil
}

// ContainerShell runs dummy shell.
func (d Dummy) ContainerShell(id string, user string, command []string, stdin io.Reader) error {
	if !d.hasContainer(id) {
		return fmt.Errorf("container %s not running", id)
	}
	return nil
}

// ContainerStatus gets dummy status.
func (d Dummy) ContainerStatus(id string) (Status, error) {
	return Status{Running: d.hasContainer(id), IPAddress: "1.1.1.1", Name: id}, nil
}

// ContainerUpload uploads to dummy container.
func (d Dummy) ContainerUpload(id string, path string, r io.Reader) error {
	if !d.hasContainer(id) {
		return fmt.Errorf("container %s not running", id)
	}
	return nil
}

// ContainerLog returns dummy logs.
func (d Dummy) ContainerLog(id string, follow bool) (io.ReadCloser, error) {
	if !d.hasContainer(id) {
		return nil, fmt.Errorf("container %s not running", id)
	}
	return ioutil.NopCloser(bytes.NewReader([]byte("hello world"))), nil
}

// ContainerCommit commits dummy container.
func (d Dummy) ContainerCommit(id string) error {
	if !d.hasContainer(id) {
		return fmt.Errorf("container %s not running", id)
	}
	return nil
}

// ContainerDeleteCommit deletes dummy commit.
func (d Dummy) ContainerDeleteCommit(id string) error {
	if d.hasContainer(id) {
		return fmt.Errorf("container %s is running", id)
	}
	return nil
}

// ImagePull pulls dummy images.
func (d Dummy) ImagePull(c []Config) error {
	return nil
}

// ProjectStop stops dummy containers.
func (d Dummy) ProjectStop(pid string) error {
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	containers := make([]string, 0)
	for _, c := range d.Tracker.Containers {
		if !strings.Contains(c, pid) {
			containers = append(containers, c)
		}
	}
	d.Tracker.Containers = containers
	return nil
}

// ProjectPurge purges dummy resources.
func (d Dummy) ProjectPurge(pid string) error {
	d.ProjectStop(pid)
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	volumes := make([]string, 0)
	for _, c := range d.Tracker.Volumes {
		if !strings.Contains(c, pid) {
			volumes = append(volumes, c)
		}
	}
	d.Tracker.Volumes = volumes
	return nil
}

// AllStop stops dummy containers.
func (d Dummy) AllStop() error {
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	d.Tracker.Containers = make([]string, 0)
	return nil
}

// AllPurge purges dummy resources.
func (d Dummy) AllPurge() error {
	d.Tracker.Sync.Lock()
	defer d.Tracker.Sync.Unlock()
	d.Tracker.Volumes = make([]string, 0)
	return nil
}
