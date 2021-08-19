package aggregate

import "github.com/volcengine/cloud-monitor-agent/monitor/model"

// AggregatorInterface for Metric data aggregate
// maybe use in the future.
type AggregatorInterface interface {
	Aggregate(metricSlice model.InputMetricSlice) *model.InputMetric
}
