package cmd

import (
	"github.com/spf13/cobra"
	"github.com/volcengine/cloud-monitor-agent/monitor/display"
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Print the metrics data to console",
	Long:  `Print the metrics data to console`,
	Run: func(cmd *cobra.Command, args []string) {
		display.CollectToConsole()
	},
}
