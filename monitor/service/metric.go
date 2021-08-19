package service

import (
	"fmt"
	"time"

	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/collectors"
	"github.com/volcengine/cloud-monitor-agent/monitor/config"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"github.com/volcengine/cloud-monitor-agent/monitor/report"
	"github.com/volcengine/cloud-monitor-agent/utils"
	"go.uber.org/zap"
)

type pluginCollectorUnloader func()

var (
	collectorList = []collectors.CollectorInterface{
		&collectors.CPUCollector{},
		&collectors.LoadCollector{},
		&collectors.MemCollector{},
		&collectors.DiskCollector{},
		&collectors.NetCollector{},
	}
	pluginCollectorUnloaders = []pluginCollectorUnloader{}
)

func loadPluginCollectors() {
	// according to the config to add plugin collector to the List
	if config.GetPluginConfig().GPUCollectorOpen {
		gpuCollector := &collectors.GPUCollector{}
		if gpuCollector.AcquireResource() {
			collectorList = append(collectorList, gpuCollector)
			pluginCollectorUnloaders = append(pluginCollectorUnloaders, gpuCollector.ReleaseResource)
		}
	}
}

func unloadPluginCollectors() {
	for _, unload := range pluginCollectorUnloaders {
		unload()
	}
}

// CollectToServer start collect task.
func CollectToServer() {
	ticker := time.NewTicker(time.Duration(config.DefaultMetricDeltaDataTimeInSecond) * time.Second)

	defer func() {
		if err := recover(); err != nil {
			logs.GetLogger().Panic("collect task panic", zap.String("reason", fmt.Sprintf("%v", err)))
		}
	}()

	loadPluginCollectors()
	defer unloadPluginCollectors()

	for range ticker.C {
		go report.SendMetricData(BuildURL(), collectMetricData())
	}
}

func collectMetricData() *model.InputMetric {
	var metricDatas []model.Metric
	now := utils.GetCurrTSInNano()
	for _, collector := range collectorList {
		collectedData := collector.Collect(now)
		if collectedData != nil {
			metricDatas = append(metricDatas, collectedData.Data...)
		}
	}

	data := &model.InputMetric{
		CollectTime: now,
		Data:        metricDatas,
	}
	return data
}

// CollectMetricDataForDisplay collect metrics data for print to console
func CollectMetricDataForDisplay() []*model.InputMetric {
	var metricDatas []*model.InputMetric
	now := utils.GetCurrTSInNano()
	for _, collector := range collectorList {
		collectedData := collector.Collect(now)
		if collectedData != nil && len(collectedData.Data) > 0 {
			metricDatas = append(metricDatas, collectedData)
		}
	}
	return metricDatas
}

// BuildURL return URL string of gateway.
func BuildURL() string {
	return config.GetMonitorConfig().Endpoint
}
