package project

import (
	"runtime"
)

func isMacOS() bool {
	return runtime.GOOS == "darwin"
}
