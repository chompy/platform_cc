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

package tests

import (
	"path"
	"strings"
	"testing"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

func TestStartStopProject(t *testing.T) {
	projectPath := path.Join("data", "sample2")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	dockerClient := NewMockDockerClient()
	p.SetDockerClient(dockerClient)
	p.Start()
	list, _ := dockerClient.GetProjectContainers(p.ID)
	assertEqual(
		strings.Contains(list[0].ID, p.Apps[0].Name),
		true,
		"expected test_app2 container",
		t,
	)
	p.Stop()
	list, _ = dockerClient.GetProjectContainers(p.ID)
	assertEqual(
		len(list),
		0,
		"expected no active containers",
		t,
	)
}

func TestStartMultipleProjects(t *testing.T) {
	// create a docker client to perform tests with
	dockerClient := NewMockDockerClient()
	// load sample1 project
	p1Path := path.Join("data", "sample1")
	p1, e := project.LoadFromPath(p1Path, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p1.SetDockerClient(dockerClient)
	p1.ID = "sample1"
	// load sample2 project
	p2Path := path.Join("data", "sample2")
	p2, e := project.LoadFromPath(p2Path, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p2.SetDockerClient(dockerClient)
	p2.ID = "sample2"
	// start projects
	p1.Start()
	p2.Start()
	// verify number of containers for sample1
	list, _ := dockerClient.GetProjectContainers(p1.ID)
	assertEqual(
		len(list),
		5,
		"unexpected number of containers for sample1 project",
		t,
	)
	// verify number of containers for sample2
	list, _ = dockerClient.GetProjectContainers(p2.ID)
	assertEqual(
		len(list),
		1,
		"unexpected number of containers for sample2 project",
		t,
	)
	// verify number of containers
	list, _ = dockerClient.GetAllContainers()
	assertEqual(
		len(list),
		6,
		"unexpected number of containers",
		t,
	)
	// verify number of volumes for sample1
	volumes, _ := dockerClient.GetProjectVolumes(p1.ID)
	assertEqual(
		len(volumes.Volumes),
		5,
		"unexpected number of volumes for sample1",
		t,
	)
	// verify number of volumes for sample2
	volumes, _ = dockerClient.GetProjectVolumes(p2.ID)
	assertEqual(
		len(volumes.Volumes),
		1,
		"unexpected number of volumes for sample2",
		t,
	)
	// stop sample1
	p1.Stop()
	// verify number of remaining containers
	list, _ = dockerClient.GetAllContainers()
	assertEqual(
		len(list),
		1,
		"unexpected number of containers",
		t,
	)
	// verify number of volumes
	volumes, _ = dockerClient.GetAllVolumes()
	assertEqual(
		len(volumes.Volumes),
		6,
		"unexpected number of volumes",
		t,
	)
	// purge projects
	p1.Purge()
	p2.Purge()
	// verify that no volumes exist
	volumes, _ = dockerClient.GetAllVolumes()
	assertEqual(
		len(volumes.Volumes),
		0,
		"unexpected number of volumes",
		t,
	)
}
