package collectors

import (
	"runtime"

	"github.com/shirou/gopsutil/load"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"go.uber.org/zap"
)

// LoadCollector is the collector type for cpu load metric.
type LoadCollector struct {
}

// Collect implement the load Collector.
func (l *LoadCollector) Collect(collectTime int64) *model.InputMetric {
	loadAvg, err := load.Avg()
	if nil != err {
		logs.GetLogger().Info("get system load error", zap.Error(err))
		return nil
	}

	numCPU := float64(runtime.NumCPU())

	metricsDatas := []model.Metric{
		{MetricName: "Load1m", MetricValue: loadAvg.Load1},
		{MetricName: "Load5m", MetricValue: loadAvg.Load5},
		{MetricName: "Load15m", MetricValue: loadAvg.Load15},
		{MetricName: "LoadPerCore1m", MetricValue: loadAvg.Load1 / numCPU},
		{MetricName: "LoadPerCore5m", MetricValue: loadAvg.Load5 / numCPU},
		{MetricName: "LoadPerCore15m", MetricValue: loadAvg.Load15 / numCPU},
	}

	return &model.InputMetric{
		Data:        metricsDatas,
		Type:        "load",
		CollectTime: collectTime,
	}
}
