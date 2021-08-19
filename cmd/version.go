package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cloud-monitor-agent",
	Long:  `Print the version number of cloud-monitor-agent`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Cloud-monitor-agent version is: %s\n", utils.AgentVersion)
	},
}
