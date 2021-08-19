package model

import (
	"io/ioutil"

	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
)

// resourceID is InstanceUUID.
var (
	resourceID string
)

func init() {
	resourceID = getResourceID()
}

const (
	// ByteToBit shift Byte to Bit
	ByteToBit = 8
	// KBToByte shift KB to Bit
	KBToByte = 1024
	// MBToByte shift MB to Byte
	MBToByte = 1024 * 1024
	// GBToByte shift GB to Byte
	GBToByte = 1024 * 1024 * 1024
	// ToPercent shift to %
	ToPercent = 100
	// WattsTomW Watts To Milli watts
	WattsTomW = 1000
)

// Metric the type for metric data.
type Metric struct {
	MetricName   string  `json:"metric_name"`
	MetricValue  float64 `json:"metric_value"`
	MetricPrefix string  `json:"metric_prefix,omitempty"`
}

// InputMetric the type for input metric.
type InputMetric struct {
	CollectTime int64    `json:"collect_time"`
	Type        string   `json:"-"`
	Data        []Metric `json:"data"`
}

// Dimension represent monitory point.
type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MetricType the type for metric.
type MetricType struct {
	// Dimensions only have one element now.
	Dimensions []Dimension `json:"dimensions"`
	MetricName string      `json:"metric_name"`
}

// MetricData the type for post metric data.
type MetricData struct {
	Metric      MetricType `json:"metric"`
	CollectTime int64      `json:"collect_time"`
	Value       float64    `json:"value"`
}

// InputMetricSlice the type for input metric slice.
type InputMetricSlice []*InputMetric

// MetricDataArr the type for metric data array.
type MetricDataArr []MetricData

// BuildMetricData will deal with InputMetric to MetricDataArr.
func BuildMetricData(inputMetric *InputMetric) MetricDataArr {
	var (
		dimension     Dimension
		metricDataArr MetricDataArr
	)

	dimension.Name = utils.UUIDName
	dimension.Value = resourceID
	if dimension.Value == "" {
		return nil
	}
	collectTime := inputMetric.CollectTime

	for _, metric := range inputMetric.Data {

		var newMetricData MetricData
		newMetricData.CollectTime = collectTime

		newMetricData.Metric.MetricName = metric.MetricName
		var dimensions []Dimension
		dimensions = append(dimensions, dimension)

		if metric.MetricPrefix != "" {
			d := Dimension{
				Name:  "DeviceName",
				Value: metric.MetricPrefix,
			}
			dimensions = append(dimensions, d)
		}
		newMetricData.Metric.Dimensions = dimensions
		// keep two decimal
		newMetricData.Value = utils.Keep2Decimal(metric.MetricValue)
		metricDataArr = append(metricDataArr, newMetricData)
	}
	return metricDataArr
}

// GetResourceId return the uuid of instance.
func getResourceID() string {
	id, err := ioutil.ReadFile(utils.ResourceIDAbsFilePath)
	if err != nil {
		logs.GetLogger().Error("the file is not exist")
		return ""
	}
	return string(id[:len(id)-1])
}
