package utils

const (
	// CronJobTimeSecond Collect the metric per 60s.
	CronJobTimeSecond = 60

	// UUIDName DimensionName is a const value.
	UUIDName = "ResourceId"

	// ResourceIDAbsFilePath This file will contain the resourceId(InstanceId).
	ResourceIDAbsFilePath = "/var/lib/cloud/data/instance-id"

	// CloudMonitorService name of influxDB which registered on gateway.
	CloudMonitorService = "x"

	// MeasurementPrefix used to distinguish out-of-ecs indicators.
	MeasurementPrefix = "Inner"

	// ConfMonitorName .
	ConfMonitorName = "conf_monitor.json"

	// ConfPluginName .
	ConfPluginName = "conf_plugin.json"
)
