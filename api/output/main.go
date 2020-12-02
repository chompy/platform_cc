package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ztrue/tracerr"
)

const levelSpacer = "   "

// Enable is a flag that sets whether or not to enable terminal output.
var Enable = false

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
	if !Enable {
		return
	}
	levelMsg(colorSuccess(msg))
}

// Warn prints a warning message to the terminal.
func Warn(msg string) {
	if !Enable {
		return
	}
	IndentLevel++
	levelMsg(colorWarn(msg))
	IndentLevel--
}

// Duration prints information and returns channel
func Duration(msg string) func() {
	if !Enable {
		return func() {}
	}
	start := time.Now()
	Info(msg)
	IndentLevel++
	done := func() {
		dur := time.Now().Sub(start)
		Info(
			fmt.Sprintf("Done (%dms).", dur.Milliseconds()),
		)
		IndentLevel--
	}
	return done
}

// Error prints an error message to the terminal and then exists.
func Error(err error) {
	if !Enable || err == nil {
		return
	}
	fmt.Println(colorError("\n\n== ERROR ==\n" + err.Error() + "\n\n"))
	tracerr.PrintSourceColor(err)
	os.Exit(1)
}
