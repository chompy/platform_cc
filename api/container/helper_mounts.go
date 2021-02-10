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

var volumeSlotRegex = regexp.MustCompile("(.*)-([0-9]{1,})")

// getMountName generates a mount name for given project id and container name.
func getMountName(pid string, name string, containerType ObjectContainerType) string {
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
