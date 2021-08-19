package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helpCmd)
}

var helpCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
