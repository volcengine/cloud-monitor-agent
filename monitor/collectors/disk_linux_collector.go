package collectors

import (
	"strings"
	"sync"

	"github.com/shirou/gopsutil/disk"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/config"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	monitorUtil "github.com/volcengine/cloud-monitor-agent/monitor/utils"
	"go.uber.org/zap"
)

// DiskCollector is the collector type for disk metric.
type DiskCollector struct {
	DiskMap sync.Map
}

// DiskIOCountersStat is the type for store disk IO data.
type DiskIOCountersStat struct {
	collectTime     int64
	uptimeInSeconds int64
	readBytes       float64
	readCount       float64
	writeBytes      float64
	writeCount      float64
}

// Collect implement the disk Collector
// we just collect the partition
// eg:sda, and it's partition is sda1,we just show the sda1.
func (d *DiskCollector) Collect(collectTime int64) *model.InputMetric {
	var (
		result       model.InputMetric
		metricsDatas []model.Metric
	)

	partitions, err := disk.Partitions(false)
	if err != nil {
		logs.GetLogger().Error("get disk partitions failed", zap.Error(err))
		return &result
	}

	diskStats, err := disk.IOCounters()
	if err != nil {
		logs.GetLogger().Warn("get disk IOCounters failed", zap.Error(err))
	}

	if err != nil {
		logs.GetLogger().Warn("get device map error", zap.Error(err))
	}

	for _, p := range partitions {
		diskMountPoint := p.Mountpoint
		// get the proper device name by mount point
		deviceName := getPartitionName(p, diskStats)
		if deviceName == "" {
			continue
		}

		// file system usage
		diskUsageMetrics := getDiskUsageMetrics(deviceName, diskMountPoint)
		metricsDatas = append(metricsDatas, diskUsageMetrics...)

		// get current stats for the device name above
		nowStats := getStats(collectTime, deviceName, diskStats)

		// calculate disk metrics
		if lastStatsData, ok := d.DiskMap.Load(deviceName); ok {
			if priorStats, ok := lastStatsData.(*DiskIOCountersStat); ok {
				diskIOMetrics := getDiskIOMetrics(deviceName, nowStats, priorStats)
				metricsDatas = append(metricsDatas, diskIOMetrics...)
			} else {
				logs.GetLogger().Error("Disk stats found in map for device" +
					", but convert failed")
			}
		} else {
			logs.GetLogger().Debug("Disk stats NOT found in map for device:" + deviceName)
		}

		d.DiskMap.Store(deviceName, nowStats)
	}

	result.Data = metricsDatas
	result.CollectTime = collectTime
	result.Type = "disk"

	return &result
}

// getDeltaTime returns the delta time between two collect.
func getDeltaTime(current, prior *DiskIOCountersStat) float64 {
	var deltaTime = float64(config.DefaultMetricDeltaDataTimeInSecond)
	collectDeltaTime := float64(current.collectTime-prior.collectTime) / 1000
	if current.uptimeInSeconds != -1 && prior.uptimeInSeconds != -1 {
		deltaTime = float64(current.uptimeInSeconds - prior.uptimeInSeconds)
	} else if collectDeltaTime > 0 {
		deltaTime = collectDeltaTime
	}

	return deltaTime
}

// getDiskUsageMetrics returns the metric for fs usage.
func getDiskUsageMetrics(partitionName, mountPoint string) []model.Metric {
	var metrics []model.Metric

	usageStats, err := disk.Usage(mountPoint)
	if err != nil {
		logs.GetLogger().Error("get disk usage failed", zap.Error(err))
		return metrics
	}

	metrics = append(metrics, model.Metric{
		MetricName:   "DiskTotal",
		MetricValue:  float64(usageStats.Total),
		MetricPrefix: partitionName,
	}, model.Metric{
		MetricName:   "DiskUsageAvail",
		MetricValue:  float64(usageStats.Free),
		MetricPrefix: partitionName,
	}, model.Metric{
		MetricName:   "DiskUsed",
		MetricValue:  float64(usageStats.Used),
		MetricPrefix: partitionName,
	}, model.Metric{
		MetricName:   "DiskUsageUtilization",
		MetricValue:  usageStats.UsedPercent,
		MetricPrefix: partitionName,
	}, model.Metric{
		MetricName:   "DiskInodesUsedPercent",
		MetricValue:  usageStats.InodesUsedPercent,
		MetricPrefix: partitionName,
	})

	return metrics
}

// getPartitionName returns the parition name ,eg: innput:/dev/sda return:sda.
func getPartitionName(p disk.PartitionStat, diskStats map[string]disk.IOCountersStat) string {
	deviceName := strings.TrimPrefix(p.Device, "/dev/")
	if _, ok := diskStats[deviceName]; ok {
		logs.GetLogger().Debug("Device name is", zap.String("deviceName", deviceName))
		return deviceName
	}
	logs.GetLogger().Warn("Device doesn't have diskStats", zap.String("Device", p.Device))

	return ""
}

// getStats returns disk stats for certain device.
func getStats(collectTime int64, deviceName string, diskStats map[string]disk.IOCountersStat) *DiskIOCountersStat {
	if _, ok := diskStats[deviceName]; ok {
		uptimeInSeconds, err := monitorUtil.GetUptimeInSeconds()

		if err != nil {
			logs.GetLogger().Error("exec GetUptimeInSeconds error", zap.Error(err))
		}

		return &DiskIOCountersStat{
			collectTime:     collectTime,
			uptimeInSeconds: uptimeInSeconds,
			readBytes:       float64(diskStats[deviceName].ReadBytes),
			readCount:       float64(diskStats[deviceName].ReadCount),
			writeBytes:      float64(diskStats[deviceName].WriteBytes),
			writeCount:      float64(diskStats[deviceName].WriteCount),
		}
	}

	return &DiskIOCountersStat{}
}

// getDiskIOMetrics return metrics by partition name.
func getDiskIOMetrics(diskPrefix string, c, l *DiskIOCountersStat) []model.Metric {
	var (
		metricsDatas    []model.Metric
		deltaReadBytes  = monitorUtil.Float64From32Bits(c.readBytes - l.readBytes)
		deltaReadReq    = monitorUtil.Float64From32Bits(c.readCount - l.readCount)
		deltaWriteBytes = monitorUtil.Float64From32Bits(c.writeBytes - l.writeBytes)
		deltaWriteReq   = monitorUtil.Float64From32Bits(c.writeCount - l.writeCount)
		// second
		deltaTime = getDeltaTime(c, l)
	)
	if deltaTime > 0 {
		metricsDatas = append(metricsDatas, model.Metric{
			MetricName:   "DiskReadBytes",
			MetricValue:  deltaReadBytes / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "DiskReadIops",
			MetricValue:  deltaReadReq / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "DiskWriteBytes",
			MetricValue:  deltaWriteBytes / deltaTime,
			MetricPrefix: diskPrefix,
		}, model.Metric{
			MetricName:   "DiskWriteIops",
			MetricValue:  deltaWriteReq / deltaTime,
			MetricPrefix: diskPrefix,
		})
	}

	return metricsDatas
}
