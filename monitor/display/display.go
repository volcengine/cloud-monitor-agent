package display

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"github.com/volcengine/cloud-monitor-agent/monitor/service"
)

// CollectToConsole TODO there need to support custom the duration.
func CollectToConsole() {
	var line int
	go disPlayMetrics(&line)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	sig := <-sigChan
	// dealwith ctrl+c signal
	if sig == syscall.SIGINT {
		cursorReset(strconv.Itoa(line + 1))
	}
}

func getMaxDataLen(datas []*model.InputMetric) int {
	var max = 0
	for _, data := range datas {
		if max < len(data.Data) {
			max = len(data.Data)
		}
	}

	return max
}

// printTitle print metadata for metrics
func printTitle(datas []*model.InputMetric) {
	var title string
	// print title.
	for _, data := range datas {
		title += fmt.Sprintf("%-*s", 33, data.Type)
	}
	fmt.Printf("%c[7;40;36m%s%c[0m\n", 0x1B, strings.ToUpper(title), 0x1B)
}

func printData(datas []*model.InputMetric) int {
	var dimensionIndex = 0
	var dataIndex = 0

	maxLen := getMaxDataLen(datas)
	// BFS display the metricsData.
	for {
		if dataIndex >= maxLen {
			break
		}

		data := datas[dimensionIndex]

		if dataIndex > len(data.Data)-1 {
			fmt.Printf("%-*s", 32, "")
			dimensionIndex++
			dimensionIndex %= len(datas)

			if dimensionIndex == 0 {
				dataIndex++
				fmt.Printf("\n")
			}

			continue
		}

		var oneLine string

		if data.Data[dataIndex].MetricPrefix != "" {
			fmt.Printf("%c[1;40;32m%s%c[0m", 0x1B, data.Data[dataIndex].MetricPrefix+":", 0x1B)
		}

		oneLine += fmt.Sprintf("%s:%.2f", data.Data[dataIndex].MetricName, data.Data[dataIndex].MetricValue)
		fmt.Printf("%-*s", 32, oneLine)

		dimensionIndex++
		dimensionIndex %= len(datas)

		if dimensionIndex == 0 {
			dataIndex++
			fmt.Printf("\n")
		}
	}

	cursorBackPriLocation(strconv.Itoa(maxLen + 1))

	return maxLen
}

func disPlayMetrics(line *int) {
	for {
		metricDatas := service.CollectMetricDataForDisplay()
		printTitle(metricDatas)
		*line = printData(metricDatas)

		time.Sleep(3 * time.Second)
	}
}
