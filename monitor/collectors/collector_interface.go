package collectors

import "github.com/volcengine/cloud-monitor-agent/monitor/model"

// CollectorInterface for raw(no aggregator) metric collect.
type CollectorInterface interface {
	Collect(collectTime int64) *model.InputMetric
}
