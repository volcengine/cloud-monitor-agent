package collectors

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"go.uber.org/zap"
)

// MemCollector is the collector type for memory metric.
type MemCollector struct {
}

// Collect implement the memory Collector.
func (m *MemCollector) Collect(collectTime int64) *model.InputMetric {
	memory, err := mem.VirtualMemory()
	if nil != err {
		logs.GetLogger().Error("Get memory stats failed", zap.Error(err))
		return nil
	}

	metricsDatas := []model.Metric{
		{
			MetricName:  "MemoryTotalSpace",
			MetricValue: float64(memory.Total),
		},
		{
			MetricName:  "MemoryFreeSpace",
			MetricValue: float64(memory.Available),
		},
		{
			MetricName:  "MemoryUsedSpace",
			MetricValue: float64(memory.Total - memory.Available),
		},
		{
			MetricName:  "MemoryBuffers",
			MetricValue: float64(memory.Buffers),
		},
		{
			MetricName:  "MemoryCached",
			MetricValue: float64(memory.Cached),
		},
		{
			MetricName:  "MemoryUsedUtilization",
			MetricValue: float64(memory.Total-memory.Available) / float64(memory.Total) * model.ToPercent,
		},
		{
			MetricName:  "MemoryFreeUtilization",
			MetricValue: float64(memory.Available) / float64(memory.Total) * model.ToPercent,
		},
	}

	return &model.InputMetric{
		Data:        metricsDatas,
		Type:        "memory",
		CollectTime: collectTime,
	}
}
