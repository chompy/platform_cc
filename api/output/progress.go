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
	"strings"
)

const progressPadChar = "."
const progressPadLength = 4

// ProgressMessageState is the state of a progress message.
type ProgressMessageState int

const (
	// ProgressMessageWait is a waiting state.
	ProgressMessageWait ProgressMessageState = 1 << iota
	// ProgressMessageDone is a done state.
	ProgressMessageDone
	// ProgressMessageSkip is a skipped state.
	ProgressMessageSkip
	// ProgressMessageError is an error state.
	ProgressMessageError
	// ProgressMessageCancel is a canceled state.
	ProgressMessageCancel
)

// String returns a string representation of the state.
func (p ProgressMessageState) String() string {
	switch p {
	case ProgressMessageWait:
		{
			return colorProgress("WAIT")
		}
	case ProgressMessageDone:
		{
			return colorSuccess("DONE")
		}
	case ProgressMessageSkip:
		{
			return colorWarn("SKIPPED")
		}
	case ProgressMessageError:
		{
			return colorError("ERROR")
		}
	case ProgressMessageCancel:
		{
			return colorWarn("CANCELED")
		}
	}
	return "???"
}

func progressMaxWidth(msgs []string) int {
	maxWidth := 0
	for _, msg := range msgs {
		msgLen := len(msg)
		if msgLen > maxWidth {
			maxWidth = msgLen
		}
	}
	return maxWidth
}

func progressPadWidth(msgs []string, index int) int {
	maxWidth := progressMaxWidth(msgs)
	padWidth := maxWidth - len(msgs[index])
	if padWidth < 0 {
		padWidth = 0
	}
	padWidth += progressPadLength
	return padWidth
}

func progressPrint(msgs []string, states []ProgressMessageState, index int) {
	msg := msgs[index]
	state := states[index]
	padWidth := progressPadWidth(msgs, index)
	// print
	levelMsg(
		msg + strings.Repeat(progressPadChar, padWidth) +
			state.String(),
	)
}

func progressPrintAll(msgs []string, states []ProgressMessageState) {
	if !isTTY() {
		return
	}
	for i := range msgs {
		progressPrint(msgs, states, i)
	}
}

func progressReprint(msgs []string, states []ProgressMessageState) {
	fmt.Printf("\033[%dA", len(msgs))
	progressPrintAll(msgs, states)
}

// Progress prints progress messages with state and returns function that updates progress when called.
func Progress(msgs []string) func(i int, s ProgressMessageState) {
	states := make([]ProgressMessageState, len(msgs))
	for i := range states {
		states[i] = ProgressMessageWait
	}
	startIndent := IndentLevel
	progressPrintAll(msgs, states)
	return func(i int, s ProgressMessageState) {
		if i < 0 || i >= len(msgs) {
			return
		}
		currentIndent := IndentLevel
		IndentLevel = startIndent
		states[i] = s
		if !isTTY() {
			progressPrint(msgs, states, i)
			return
		}
		progressReprint(msgs, states)
		IndentLevel = currentIndent
	}
}
