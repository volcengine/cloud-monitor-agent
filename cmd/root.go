package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloud-monitor-agent",
	Short: "Cloud-monitor-agent is a collector working guest OS",
	Long:  `Cloud-monitor-agent is a collector working guest OS`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute is the exe of entrance
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
