package main

import (
	"gitlab.com/contextualcode/platform_cc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
