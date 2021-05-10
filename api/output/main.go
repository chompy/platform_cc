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

	"golang.org/x/term"
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

// IsTTY returns true if running with a TTY.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// WriteStdout writes a string to STDOUT.
func WriteStdout(msg string) {
	os.Stdout.WriteString(msg)
}

// WriteStderr writes a string to STDERR.
func WriteStderr(msg string) {
	os.Stderr.WriteString(msg)
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
	WriteStderr(
		strings.Repeat(levelSpacer, level) + msg + "\n",
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
			dur := time.Since(start)
			LogInfo(msg + fmt.Sprintf(" (%dms).", dur.Milliseconds()))
		}
	}
	levelMsg(msg)
	IndentLevel++
	done := func() {
		dur := time.Since(start)
		levelMsg(
			colorSuccess(fmt.Sprintf("Done (%dms).", dur.Milliseconds())),
		)
		LogInfo(msg + fmt.Sprintf(" (%dms).", dur.Milliseconds()))
		IndentLevel--
	}
	return done
}

// ContainerLog prints container log line to stdout.
func ContainerLog(name string, msg string) {
	WriteStdout(colorSuccess(fmt.Sprintf("[%s] ", name)) + msg + "\n")
}
