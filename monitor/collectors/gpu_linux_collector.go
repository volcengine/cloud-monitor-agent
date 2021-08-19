//+build linux

package collectors

import (
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"go.uber.org/zap"
)

const PREFIX = "gpu_"

// GPUCollector is the collector type for gpu metric.
type GPUCollector struct {
}

func (GPUCollector) AcquireResource() bool {
	if ret := nvml.Init(); ret != nvml.SUCCESS {
		logs.GetLogger().Error("failed to init nvml")
		return false
	}
	return true
}

func (GPUCollector) ReleaseResource() {
	nvml.Shutdown()
}

// Collect implement the GPU Collector.
func (c *GPUCollector) Collect(collectTime int64) *model.InputMetric {
	var (
		result       model.InputMetric
		metricsDatas []model.Metric
	)

	deviceCount, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		logs.GetLogger().Error("get gpu device count error: ", zap.Int("nvml return value", int(ret)))
		return nil
	}
	if deviceCount == 0 {
		logs.GetLogger().Warn("no gpu bound on nvidia driver")
		return nil
	}

	for i := 0; i < deviceCount; i++ {
		handle, _ := nvml.DeviceGetHandleByIndex(i)
		decoderUtil, _, _ := nvml.DeviceGetDecoderUtilization(handle)
		encoderUtil, _, _ := nvml.DeviceGetEncoderUtilization(handle)
		temperature, _ := nvml.DeviceGetTemperature(handle, nvml.TEMPERATURE_GPU)
		gpuUsedUtil, _ := nvml.DeviceGetUtilizationRates(handle)
		memInfo, _ := nvml.DeviceGetMemoryInfo(handle)
		memFree := float64(memInfo.Free)
		memTotal := float64(memInfo.Total)
		memUsed := float64(memInfo.Used)
		pwr, _ := nvml.DeviceGetPowerUsage(handle) // power draw (unit32, in milliwatts)
		gpuID := PREFIX + strconv.Itoa(i)

		metricsInfo := []model.Metric{
			{MetricName: "GpuDecoderUtilization", MetricValue: float64(decoderUtil), MetricPrefix: gpuID},
			{MetricName: "GpuEncoderUtilization", MetricValue: float64(encoderUtil), MetricPrefix: gpuID},
			{MetricName: "GpuTemperature", MetricValue: float64(temperature), MetricPrefix: gpuID},
			{MetricName: "GpuUsedUtilization", MetricValue: float64(gpuUsedUtil.Gpu), MetricPrefix: gpuID},
			{MetricName: "GpuMemoryFreeSpace", MetricValue: memFree / model.MBToByte, MetricPrefix: gpuID},
			{MetricName: "GpuMemoryFreeUtilization", MetricValue: model.ToPercent * (memFree / memTotal), MetricPrefix: gpuID},
			{MetricName: "GpuMemoryTotalSpace", MetricValue: memTotal, MetricPrefix: gpuID},
			{MetricName: "GpuMemoryUsedSpace", MetricValue: memUsed, MetricPrefix: gpuID},
			{MetricName: "GpuMemoryUsedUtilization", MetricValue: model.ToPercent * (memUsed / memTotal), MetricPrefix: gpuID},
			{MetricName: "GpuPowerReadingsPowerDraw", MetricValue: float64(pwr) / model.WattsTomW, MetricPrefix: gpuID},
		}

		metricsDatas = append(metricsDatas, metricsInfo...)

	}

	result.Data = metricsDatas
	result.Type = "gpu"
	result.CollectTime = collectTime

	return &result
}
