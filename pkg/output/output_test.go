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

package output

import (
	"testing"
)

// TestOutputProgress tests output progress messages.
func TestOutputProgress(t *testing.T) {
	msgs := []string{
		"Test A",
		"Test B",
		"Test C",
	}
	Enable = true
	prog := Progress(msgs)
	prog(1, ProgressMessageDone, nil, nil)
	prog(2, ProgressMessageError, nil, nil)
}