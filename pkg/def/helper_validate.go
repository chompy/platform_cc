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

package def

import (
	"fmt"
	"strings"
)

func sliceContainsString(slice []string, s string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}
	return false
}

func validateMustContainOne(slice []string, s string, key string) error {
	if !sliceContainsString(slice, s) {
		return NewValidateError(
			key,
			fmt.Sprintf("must be one of: %s got: %s", strings.Join(slice, ", "), s),
		)
	}
	return nil
}

func sliceContainsInt(slice []int, n int) bool {
	for _, val := range slice {
		if val == n {
			return true
		}
	}
	return false
}

func validateMustContainOneInt(slice []int, n int, key string) error {
	if !sliceContainsInt(slice, n) {
		outList := ""
		for _, val := range slice {
			outList += fmt.Sprintf("%d, ", val)
		}
		return NewValidateError(
			key,
			fmt.Sprintf("must be one of: %s got: %d", strings.TrimRight(outList, ","), n),
		)
	}
	return nil
}
