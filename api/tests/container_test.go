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
	"fmt"
	"path"
	"strings"
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
		len(ch.Tracker.Containers),
		len(p.Apps)+len(p.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(ch.Tracker.Volumes),
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
		len(ch.Tracker.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(ch.Tracker.Volumes),
		len(p.Apps)+len(p.Services),
		"wrong number of volumes",
		t,
	)
	p.Purge()
	// ensure no volumes remain
	assertEqual(
		len(ch.Tracker.Volumes),
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
		len(ch.Tracker.Containers),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(ch.Tracker.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of volumes",
		t,
	)
	// stop p1
	p1.Stop()

	// number of "running" containers should be len apps + len services
	assertEqual(
		len(ch.Tracker.Containers),
		len(p2.Apps)+len(p2.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	assertEqual(
		len(ch.Tracker.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of volumes",
		t,
	)
	// purge and ensure all containers and volumes are deleted
	p1.Purge()
	p2.Purge()
	assertEqual(
		len(ch.Tracker.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	assertEqual(
		len(ch.Tracker.Volumes),
		0,
		"wrong number of volumes",
		t,
	)
	// test all stop
	p1.Start()
	p2.Start()
	ch.AllStop()
	assertEqual(
		len(ch.Tracker.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	assertEqual(
		len(ch.Tracker.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of volumes",
		t,
	)
	// test all purge
	ch.AllPurge()
	assertEqual(
		len(ch.Tracker.Volumes),
		0,
		"wrong number of volumes",
		t,
	)
}

func TestStartProjectSlots(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	pPath := path.Join("data", "sample1")
	p, e := project.LoadFromPath(pPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p.SetContainerHandler(ch)
	p.ID = "sample1"
	// set slot 2 and try purging
	p.SetSlot(2)
	p.Start()
	assertEqual(
		len(ch.Tracker.Volumes),
		5,
		"expected five container volumes",
		t,
	)
	assertEqual(
		strings.Contains(ch.Tracker.Volumes[0], "-2"),
		true,
		"expected volume to contain slot suffix (-2)",
		t,
	)
	p.PurgeSlot()
	assertEqual(
		len(ch.Tracker.Volumes),
		0,
		"expected zero container volumes",
		t,
	)
	// start project twice with slot 1 and slot 2
	p.Start()
	p.Stop()
	p.SetSlot(1)
	p.Start()
	p.Stop()
	// try to purge slot 1, expect error
	e = p.PurgeSlot()
	assertEqual(
		e != nil && strings.Contains(e.Error(), "cannot delete"),
		true,
		"expected cannot delete slot error",
		t,
	)
	// ensure number of volumes
	assertEqual(
		len(ch.Tracker.Volumes),
		10,
		"expected 10 volumes",
		t,
	)
	// delete slot 2
	p.SetSlot(2)
	p.PurgeSlot()
	assertEqual(
		len(ch.Tracker.Volumes),
		5,
		"expected five volumes",
		t,
	)
}

func TestSlotCopy(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	pPath := path.Join("data", "sample1")
	p, e := project.LoadFromPath(pPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p.SetContainerHandler(ch)
	p.ID = "sample1"
	// start under slot 1
	p.SetSlot(1)
	p.Start()
	// record number of volumes
	volBefore := len(ch.Tracker.Volumes)
	// perform copy
	destVol := 3
	p.CopySlot(destVol)
	// assert that number of volumes has doubled
	assertEqual(
		len(ch.Tracker.Volumes),
		volBefore*2,
		"expected number of volumes to double",
		t,
	)
	// assert that new volumes have correct slot suffix
	newVolCount := 0
	for _, volName := range ch.Tracker.Volumes {
		if strings.Contains(volName, fmt.Sprintf("-%d", destVol)) {
			newVolCount++
		}
	}
	assertEqual(
		newVolCount,
		volBefore,
		"expected new volumes to has correct slot suffix",
		t,
	)
}

func TestProjectStatus(t *testing.T) {
	projectPath := path.Join("data", "sample1")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	ch := container.NewDummy()
	p.SetContainerHandler(ch)
	p.Start()

	count := len(p.Services)
	for _, app := range p.Apps {
		count += 1 + len(app.Workers)
	}
	// project status
	s := p.Status()
	assertEqual(
		len(s),
		count,
		"unexpected number of containers returned in status",
		t,
	)
	assertEqual(
		s[0].ProjectID,
		p.ID,
		"status returned unexpected project id",
		t,
	)
	// all status
	as, _ := ch.AllStatus()
	assertEqual(
		len(as),
		count,
		"unexpected number of containers returned in all status",
		t,
	)
	assertEqual(
		as[0].ProjectID,
		p.ID,
		"all status returned unexpected project id",
		t,
	)
}

func TestAllStatus(t *testing.T) {
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
	s, _ := ch.AllStatus()
	assertEqual(
		len(s),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"unexpected number of results in all status",
		t,
	)
	for _, st := range s {
		if st.ProjectID != p1.ID && st.ProjectID != p2.ID {
			t.Error("status did not have known project id")
		}
	}
	// TODO more checks?
}