package collectors

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"go.uber.org/zap"
)

// CPUStates is the type for store cpu state
type CPUStates struct {
	user         float64
	guest        float64
	system       float64
	idle         float64
	other        float64
	nice         float64
	iowait       float64
	irq          float64
	softirq      float64
	steal        float64
	totalCPUTime float64
}

// CPUCollector is the collector type for cpu metric.
type CPUCollector struct {
	LastStates *CPUStates
}

// getTotalCPUTime return the total time of system.
func getTotalCPUTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Idle
	return total
}

// Collect implement the cpu Collector.
func (c *CPUCollector) Collect(collectTime int64) *model.InputMetric {
	cpuStats, err := cpu.Times(false)
	if err != nil || len(cpuStats) == 0 {
		logs.GetLogger().Error("get cpu stat error", zap.Error(err))
		return nil
	}

	stat := cpuStats[0]
	nowStates := &CPUStates{
		user:         stat.User,
		guest:        stat.Guest,
		system:       stat.System,
		idle:         stat.Idle,
		nice:         stat.Nice,
		iowait:       stat.Iowait,
		irq:          stat.Irq,
		softirq:      stat.Softirq,
		steal:        stat.Steal,
		totalCPUTime: getTotalCPUTime(stat),
	}

	if c.LastStates == nil {
		c.LastStates = nowStates
		return nil
	}

	totalCPUTime := getTotalCPUTime(stat)
	totalDelta := totalCPUTime - c.LastStates.totalCPUTime

	// all value is the percentage
	cpuUsageUser := model.ToPercent * (nowStates.user - c.LastStates.user -
		(nowStates.guest - c.LastStates.guest)) / totalDelta
	cpuUsageSystem := model.ToPercent * (nowStates.system - c.LastStates.system) / totalDelta
	cpuUsageIdle := model.ToPercent * (nowStates.idle - c.LastStates.idle) / totalDelta
	cpuUsageNice := model.ToPercent * (nowStates.nice - c.LastStates.nice) / totalDelta
	cpuUsageIOWait := model.ToPercent * (nowStates.iowait - c.LastStates.iowait) / totalDelta
	cpuUsageIrq := model.ToPercent * (nowStates.irq - c.LastStates.irq) / totalDelta
	cpuUsageSoftIrq := model.ToPercent * (nowStates.softirq - c.LastStates.softirq) / totalDelta
	cpuUsageSteal := model.ToPercent * (nowStates.steal - c.LastStates.steal) / totalDelta
	cpuOther := cpuUsageNice + cpuUsageSoftIrq + cpuUsageIrq + cpuUsageSteal

	c.LastStates = nowStates

	metricsDatas := []model.Metric{
		{MetricName: "CpuUser", MetricValue: cpuUsageUser},
		{MetricName: "CpuSystem", MetricValue: cpuUsageSystem},
		{MetricName: "CpuIdle", MetricValue: cpuUsageIdle},
		{MetricName: "CpuNice", MetricValue: cpuUsageNice},
		{MetricName: "CpuWait", MetricValue: cpuUsageIOWait},
		{MetricName: "CpuIrq", MetricValue: cpuUsageIrq},
		{MetricName: "CpuSoftIrq", MetricValue: cpuUsageSoftIrq},
		{MetricName: "CpuTotal", MetricValue: 100 - cpuUsageIdle},
		{MetricName: "CpuOther", MetricValue: cpuOther},
	}

	return &model.InputMetric{
		Data:        metricsDatas,
		Type:        "cpu",
		CollectTime: collectTime,
	}
}
