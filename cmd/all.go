package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api/docker"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Manage all projects.",
}

var allStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all projects.",
	Run: func(cmd *cobra.Command, args []string) {
		docker, err := docker.NewClient()
		if err != nil {
			handleError(err)
		}
		handleError(docker.DeleteAllContainers())
	},
}

var allPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all projects.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("!!! WARNING: PURGING ALL PLATFORM.CC DATA IN 5 SECONDS !!!")
		time.Sleep(5 * time.Second)
		docker, err := docker.NewClient()
		handleError(err)
		handleError(docker.DeleteAllContainers())
		handleError(docker.DeleteAllVolumes())
		handleError(docker.DeleteNetwork())
	},
}

func init() {
	allCmd.AddCommand(allStopCmd)
	allCmd.AddCommand(allPurgeCmd)
	rootCmd.AddCommand(allCmd)
}
