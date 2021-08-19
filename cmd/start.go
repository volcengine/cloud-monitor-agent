package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/manager"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cloud-monitor-agent",
	Long:  `Start cloud-monitor-agent`,
	Run: func(cmd *cobra.Command, args []string) {
		err := manager.SingleInstanceCheck()

		if err != nil {
			fmt.Println("Start failed.cloudMonitorAgent was started in before.")
			logs.GetLogger().Warn("Start failed: cloudMonitorAgent was started in before.",
				zap.Int("PID", os.Getpid()))
			return
		}

		sManager := manager.NewServiceManager()
		sManager.RegisterService()
		sManager.InitService()

		//TODO HB start to send metrics to server
		sManager.StartService()
		fmt.Println("Cloud-Monitor-Agent runs successfully.")
		logs.GetLogger().Info("cloudMonitorAgent runs successfully")
		utils.StartDaemon()
	},
}
