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
	"testing"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

func TestStartStopProject(t *testing.T) {
	projectPath := path.Join("data", "sample2")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	ch := container.NewDummy()
	p.SetContainerHandler(ch)
	p.Start()
	// number of "running" containers should be len apps + len services
	assertEqual(
		len(*ch.Containers),
		len(p.Apps)+len(p.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(*ch.Volumes),
		len(p.Apps)+len(p.Services),
		"wrong number of volumes",
		t,
	)
	// ensure first app is "running"
	s, _ := ch.ContainerStatus(p.NewContainer(p.Apps[0]).Config.GetContainerName())
	assertEqual(
		s.Running,
		true,
		"expected test_app2 container",
		t,
	)
	// stop project and ensure first app is no longer running
	p.Stop()
	s, _ = ch.ContainerStatus(p.NewContainer(p.Apps[0]).Config.GetContainerName())
	assertEqual(
		s.Running,
		false,
		"expected test_app2 container stopped",
		t,
	)
	// ensure no containers are running
	assertEqual(
		len(*ch.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(*ch.Volumes),
		len(p.Apps)+len(p.Services),
		"wrong number of volumes",
		t,
	)
	p.Purge()
	// ensure no volumes remain
	assertEqual(
		len(*ch.Volumes),
		0,
		"wrong number of volumes",
		t,
	)
}

func TestStartMultipleProjects(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	p1Path := path.Join("data", "sample1")
	p1, e := project.LoadFromPath(p1Path, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p1.SetContainerHandler(ch)
	p1.ID = "sample1"
	// load sample2 project
	p2Path := path.Join("data", "sample2")
	p2, e := project.LoadFromPath(p2Path, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p2.SetContainerHandler(ch)
	p2.ID = "sample2"
	// start projects
	p1.Start()
	p2.Start()
	// number of "running" containers should be len apps + len services
	assertEqual(
		len(*ch.Containers),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(*ch.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of volumes",
		t,
	)
	// stop p1
	p1.Stop()

	// number of "running" containers should be len apps + len services
	assertEqual(
		len(*ch.Containers),
		len(p2.Apps)+len(p2.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(*ch.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of volumes",
		t,
	)
	// purge and ensure all containers and volumes are deleted
	p1.Purge()
	p2.Purge()
	assertEqual(
		len(*ch.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	assertEqual(
		len(*ch.Volumes),
		0,
		"wrong number of volumes",
		t,
	)
	// test all stop
	p1.Start()
	p2.Start()
	ch.AllStop()
	assertEqual(
		len(*ch.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	assertEqual(
		len(*ch.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of volumes",
		t,
	)
	// test all purge
	ch.AllPurge()
	assertEqual(
		len(*ch.Volumes),
		0,
		"wrong number of volumes",
		t,
	)
}
