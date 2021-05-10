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
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func errorStackTrace(err error) string {
	out := ""
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			out += fmt.Sprintf("%+s:%d\n", f, f)
		}
	}
	return out
}

// Error prints an error message to the terminal and then exists.
func Error(err error) {
	LogError(err)
	if err == nil {
		return
	}
	WriteStdout(colorError("\nERROR:\n" + err.Error() + "\n"))
	if !Verbose {
		os.Exit(1)
	}
	WriteStdout("\n" + Color(errorStackTrace(err), 36) + "\n")
	os.Exit(1)
}

// ErrorText prints message using the error text color.
func ErrorText(msg string) {
	if !Enable {
		return
	}
	levelMsg(colorError(msg))
}
