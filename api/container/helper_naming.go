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
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// formatings
const containerNamingPrefix = "pcc-%s-"
const containerVolumeNameFormat = containerNamingPrefix + "v-%s"
const containerNameFormat = containerNamingPrefix + "%s-%s"
const containerNetworkNameFormat = containerNamingPrefix + "n"

// regex
var containerNameRegex = regexp.MustCompile("pcc-(.*)-(.*)-(.*)")
var volumeSlotRegex = regexp.MustCompile("(.*)-([0-9]{1,})")
var imageRegex = regexp.MustCompile("docker\\.registry\\.platform\\.sh\\/(.*)-(.*)")

// ObjectContainerType defines the type of container.
type ObjectContainerType byte

const (
	// ObjectContainerNone is an unknown container.
	ObjectContainerNone ObjectContainerType = '-'
	// ObjectContainerApp is an application container.
	ObjectContainerApp ObjectContainerType = 'a'
	// ObjectContainerWorker is a worker container.
	ObjectContainerWorker ObjectContainerType = 'w'
	// ObjectContainerService is a service container.
	ObjectContainerService ObjectContainerType = 's'
	// ObjectContainerRouter is the router container.
	ObjectContainerRouter ObjectContainerType = 'r'
)

// TypeName gets the type of container as a string.
func (o ObjectContainerType) TypeName() string {
	switch o {
	case ObjectContainerApp:
		{
			return "app"
		}
	case ObjectContainerWorker:
		{
			return "worker"
		}
	case ObjectContainerService:
		{
			return "service"
		}
	case ObjectContainerRouter:
		{
			return "router"
		}
	}
	return "unknown"
}

func containerName(projectID string, objectType ObjectContainerType, name string) string {
	if objectType == ObjectContainerRouter {
		return "pcc-router-1"
	}
	return fmt.Sprintf(containerNameFormat, projectID, string(objectType), name)
}

func containerConfigFromName(name string) Config {
	r := containerNameRegex.FindStringSubmatch(name)
	if len(r) < 4 {
		return Config{}
	}
	return Config{
		ProjectID:  r[1],
		ObjectType: ObjectContainerType(r[2][0]),
		ObjectName: r[3],
	}
}

// getMountName generates a mount name for given project id and container name.
func getMountName(pid string, name string, containerType ObjectContainerType) string {
	// name prefixed with underscore should be a global volume
	if name != "" && name[0] == '_' {
		return strings.TrimRight(fmt.Sprintf(containerNamingPrefix, name[1:]), "-")
	}
	return fmt.Sprintf(containerNamingPrefix+"%s-%s", pid, string(containerType), name)
}

// volumeBelongsToSlot return true if given volume belongs to given slot.
func volumeBelongsToSlot(name string, slot int) bool {
	if slot > 1 && strings.HasSuffix(name, fmt.Sprintf("-%d", slot)) {
		return true
	} else if slot <= 1 {
		return !volumeSlotRegex.MatchString(name)
	}
	return false
}

// volumeStripSlot strips the slot from the volume name.
func volumeStripSlot(name string) string {
	return volumeSlotRegex.ReplaceAllString(name, "$1")
}

// volumeWithSlot return name of volume with given slot.
func volumeWithSlot(name string, slot int) string {
	if slot <= 1 {
		return volumeStripSlot(name)
	}
	return fmt.Sprintf("%s-%d", volumeStripSlot(name), slot)
}

// volumeGetSlot returns the slot number of the given volume.
func volumeGetSlot(name string) int {
	res := volumeSlotRegex.FindStringSubmatch(name)
	if res == nil || len(res) < 3 || res[2] == "" {
		return 1
	}
	out, _ := strconv.Atoi(res[2])
	return out
}

// typeFromImageName returns the service type from the image name.
func typeFromImageName(name string) string {
	m := imageRegex.FindStringSubmatch(name)
	if m == nil || len(m) < 3 {
		return ""
	}
	return fmt.Sprintf("%s:%s", m[1], strings.Split(m[2], ":")[0])
}
