package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "platform_cc",
	Short: "Platform.cc is a tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.",
}

// Execute - run root command
func Execute() error {
	// hack that allows old style semicolon (:) seperated
	// subcommands to work
	args := make([]string, 1)
	args[0] = os.Args[0]
	args = append(args, strings.Split(os.Args[1], ":")...)
	args = append(args, os.Args[2:]...)
	os.Args = args
	return rootCmd.Execute()
}

func er(msg interface{}) {
	log.Fatal(msg)
}

func init() {

}
