package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/utils"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop  cloud-monitor-agent",
	Long:  `Stop  cloud-monitor-agent`,
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := utils.GetMainProcessPid()
		if err != nil {
			fmt.Println("Stop cloud-monitor-agent failed")
			logs.GetLogger().Error("Get main process pid failed", zap.Error(err))
			os.Exit(1)
		}

		// On Unix systems, FindProcess always succeeds and returns a Process
		// for the given pid, regardless of whether the process exists.
		process, _ := os.FindProcess(pid)

		// if stop by SIGTERM signal not work, will be try three times.
		var retryTime int = 3
		for {
			err := process.Signal(syscall.SIGTERM)
			if err != nil {
				retryTime--
				if retryTime < 0 {
					break
				}
				time.Sleep(1 * time.Second)

				continue
			}
			os.Exit(0)
		}

		// When the process is executed here
		// it means that the send SIGTERM signal has failed for three times
		// so we send the SIGKILL signal directly.
		if retryTime < 0 {
			err := process.Signal(syscall.SIGKILL)
			if err != nil {
				fmt.Println("Stop cloud-monitor-agent failed")
				logs.GetLogger().Error("SIGKILL not work, stop agent failed", zap.Error(err))
				os.Exit(1)
			}
		}
	},
}
