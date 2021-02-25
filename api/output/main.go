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

// Package output provides functions for displaying information to the end user.
package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ztrue/tracerr"
)

const levelSpacer = "   "

// Enable is a flag that enables terminal output.
var Enable = false

// Verbose is a flag that enables more verbose output.
var Verbose = false

// Logging is a flag that enables writting to log file.
var Logging = false

// IndentLevel is the current indentation level.
var IndentLevel = 0

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func levelMsg(msg string) {
	level := IndentLevel
	if level < 0 {
		level = 0
	}
	switch level {
	case 0:
		{
			msg = "> " + msg
			break
		}
	default:
		{
			msg = "- " + msg
			break
		}
	}
	os.Stderr.Write(
		[]byte(strings.Repeat(levelSpacer, level) + msg + "\n"),
	)
}

// Info prints information to the terminal.
func Info(msg string) {
	LogInfo(msg)
	if !Enable {
		return
	}
	levelMsg(colorSuccess(msg))
}

// Warn prints a warning message to the terminal.
func Warn(msg string) {
	LogWarn(msg)
	if !Enable {
		return
	}
	IndentLevel++
	levelMsg(colorWarn(msg))
	IndentLevel--
}

// Duration prints information and returns channel
func Duration(msg string) func() {
	start := time.Now()
	if !Enable {
		return func() {
			dur := time.Now().Sub(start)
			LogInfo(msg + fmt.Sprintf(" (%dms).", dur.Milliseconds()))
		}
	}
	levelMsg(msg)
	IndentLevel++
	done := func() {
		dur := time.Now().Sub(start)
		levelMsg(
			colorSuccess(fmt.Sprintf("Done (%dms).", dur.Milliseconds())),
		)
		LogInfo(msg + fmt.Sprintf(" (%dms).", dur.Milliseconds()))
		IndentLevel--
	}
	return done
}

// Error prints an error message to the terminal and then exists.
func Error(err error) {
	LogError(err)
	if err == nil {
		return
	}
	if !Verbose {
		fmt.Println(colorError("\nERROR:\n" + err.Error() + "\n"))
		os.Exit(1)
	}
	fmt.Println(colorError("\n=== ERROR ===\n"))
	// is tty
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		tracerr.PrintSourceColor(err)
		os.Exit(1)
	}
	tracerr.PrintSource(err)
	os.Exit(1)
}

// ErrorText prints message using the error text color.
func ErrorText(msg string) {
	if !Enable {
		return
	}
	levelMsg(colorError(msg))
}
