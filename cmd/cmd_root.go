package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "platform_cc",
	Short: "Platform.cc is a tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute - run root command
func Execute() error {
	return rootCmd.Execute()
}

func er(msg interface{}) {
	log.Fatal(msg)
}

func init() {

}
