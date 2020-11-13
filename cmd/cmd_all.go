package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/contextualcode/platform_cc/api"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Subcommand for managing all projects.",
}

var allStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all projects.",
	RunE: func(cmd *cobra.Command, args []string) error {
		docker, err := api.NewDockerClient()
		if err != nil {
			return err
		}
		if err := docker.DeleteAllContainers(); err != nil {
			return err
		}
		return docker.DeleteAllNetworks()
	},
}

var allPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all projects.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("!!! WARNING: PURGING ALL PLATFORM.CC DATA IN 5 SECONDS !!!")
		time.Sleep(5 * time.Second)
		docker, err := api.NewDockerClient()
		if err != nil {
			return err
		}
		if err := docker.DeleteAllContainers(); err != nil {
			return err
		}
		if err := docker.DeleteAllNetworks(); err != nil {
			return err
		}
		return docker.DeleteAllVolumes()
	},
}

func init() {
	allCmd.AddCommand(allStopCmd)
	allCmd.AddCommand(allPurgeCmd)
	rootCmd.AddCommand(allCmd)
}
