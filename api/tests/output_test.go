package tests

import (
	"testing"

	"gitlab.com/contextualcode/platform_cc/api/output"
)

// TestOutputProgress tests output progress messages.
func TestOutputProgress(t *testing.T) {

	msgs := []string{
		"Test A",
		"Test B",
		"Test C",
	}
	output.Enable = true

	prog := output.Progress(msgs)
	prog(1, output.ProgressMessageDone)
	prog(2, output.ProgressMessageError)

}
