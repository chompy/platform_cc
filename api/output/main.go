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
	fmt.Println(
		strings.Repeat(levelSpacer, level) + msg,
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
	if !Enable || err == nil {
		return
	}
	fmt.Println(colorError("\n\n== ERROR ==\n" + err.Error() + "\n\n"))
	tracerr.PrintSourceColor(err)
	os.Exit(1)
}
