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

package project

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/def"
)

func TestStartStopProject(t *testing.T) {
	projectPath := path.Join("_test_data", "sample2")
	p, e := LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	ch := container.NewDummy()
	p.SetContainerHandler(ch)
	p.Start()
	// number of "running" containers should be len apps + len services
	def.AssertEqual(
		len(ch.Tracker.Containers),
		len(p.Apps)+len(p.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		len(p.Apps)+len(p.Services)+1,
		"wrong number of volumes",
		t,
	)
	// ensure first app is "running"
	s, _ := ch.ContainerStatus(p.NewContainer(p.Apps[0]).Config.GetContainerName())
	def.AssertEqual(
		s.Running,
		true,
		"expected test_app2 container",
		t,
	)
	// stop project and ensure first app is no longer running
	p.Stop()
	s, _ = ch.ContainerStatus(p.NewContainer(p.Apps[0]).Config.GetContainerName())
	def.AssertEqual(
		s.Running,
		false,
		"expected test_app2 container stopped",
		t,
	)
	// ensure no containers are running
	def.AssertEqual(
		len(ch.Tracker.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		len(p.Apps)+len(p.Services)+1,
		"wrong number of volumes",
		t,
	)
	// ensure only global volume remains after project purge
	p.Purge()
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		1,
		"wrong number of volumes",
		t,
	)
	// ensure no volume remains after global purge
	ch.AllPurge()
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		0,
		"wrong number of volumes",
		t,
	)

}

func TestStartMultipleProjects(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	p1Path := path.Join("_test_data", "sample1")
	p1, e := LoadFromPath(p1Path, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p1.SetContainerHandler(ch)
	p1.ID = "sample1"
	// load sample2 project
	p2Path := path.Join("_test_data", "sample2")
	p2, e := LoadFromPath(p2Path, true)
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
	def.AssertEqual(
		len(ch.Tracker.Containers),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services)+1,
		"wrong number of volumes",
		t,
	)
	// stop p1
	p1.Stop()

	// number of "running" containers should be len apps + len services
	def.AssertEqual(
		len(ch.Tracker.Containers),
		len(p2.Apps)+len(p2.Services),
		"wrong number of running containers",
		t,
	)
	// check number of volumes
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services)+1,
		"wrong number of volumes",
		t,
	)
	// purge and ensure all containers and volumes are deleted
	p1.Purge()
	p2.Purge()
	def.AssertEqual(
		len(ch.Tracker.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		1,
		"wrong number of volumes",
		t,
	)
	// test all stop
	p1.Start()
	p2.Start()
	ch.AllStop()
	def.AssertEqual(
		len(ch.Tracker.Containers),
		0,
		"wrong number of running containers",
		t,
	)
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		len(p1.Apps)+len(p1.Services)+len(p2.Apps)+len(p2.Services)+1,
		"wrong number of volumes",
		t,
	)
	// test all purge
	ch.AllPurge()
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		0,
		"wrong number of volumes",
		t,
	)
}

func TestStartProjectSlots(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	pPath := path.Join("_test_data", "sample1")
	p, e := LoadFromPath(pPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p.SetContainerHandler(ch)
	p.ID = "sample1"
	// set slot 2 and try purging
	p.SetSlot(2)
	p.Start()
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		6,
		"expected 5 container volumes + 1 global",
		t,
	)
	for _, volName := range ch.Tracker.Volumes {
		if strings.Count(volName, "-") < 2 {
			continue
		}
		def.AssertEqual(
			strings.HasSuffix(volName, "-2"),
			true,
			"expected volume to contain slot suffix (-2)",
			t,
		)
	}
	p.PurgeSlot()
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		1,
		"expected one (global) container volumes",
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
	def.AssertEqual(
		e != nil && strings.Contains(e.Error(), "cannot delete"),
		true,
		"expected cannot delete slot error",
		t,
	)
	// ensure number of volumes
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		11,
		"expected 10 volumes + 1 global",
		t,
	)
	// delete slot 2
	p.SetSlot(2)
	p.PurgeSlot()
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		6,
		"expected five volumes + one global",
		t,
	)
}

func TestSlotCopy(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	pPath := path.Join("_test_data", "sample1")
	p, e := LoadFromPath(pPath, true)
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
	volBefore := len(ch.Tracker.Volumes) - 1
	// perform copy
	destVol := 3
	p.CopySlot(destVol)
	// assert that number of volumes has doubled
	def.AssertEqual(
		len(ch.Tracker.Volumes),
		(volBefore*2)+1,
		"expected number of volumes to double + global",
		t,
	)
	// assert that new volumes have correct slot suffix
	newVolCount := 0
	for _, volName := range ch.Tracker.Volumes {
		if strings.Contains(volName, fmt.Sprintf("-%d", destVol)) {
			newVolCount++
		}
	}
	def.AssertEqual(
		newVolCount,
		volBefore,
		"expected new volumes to has correct slot suffix",
		t,
	)
}

func TestProjectStatus(t *testing.T) {
	projectPath := path.Join("_test_data", "sample1")
	p, e := LoadFromPath(projectPath, true)
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
	def.AssertEqual(
		len(s),
		count,
		"unexpected number of containers returned in status",
		t,
	)
	def.AssertEqual(
		s[0].ProjectID,
		p.ID,
		"status returned unexpected project id",
		t,
	)
	// all status
	as, _ := ch.AllStatus()
	def.AssertEqual(
		len(as),
		count,
		"unexpected number of containers returned in all status",
		t,
	)
	def.AssertEqual(
		as[0].ProjectID,
		p.ID,
		"all status returned unexpected project id",
		t,
	)
}

func TestAllStatus(t *testing.T) {
	ch := container.NewDummy()
	// load sample1 project
	p1Path := path.Join("_test_data", "sample1")
	p1, e := LoadFromPath(p1Path, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	p1.SetContainerHandler(ch)
	p1.ID = "sample1"
	// load sample2 project
	p2Path := path.Join("_test_data", "sample2")
	p2, e := LoadFromPath(p2Path, true)
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
	def.AssertEqual(
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
