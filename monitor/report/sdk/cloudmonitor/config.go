package cloudmonitor

import (
	"net/http"
	"net/url"
	"time"

	"github.com/volcengine/cloud-monitor-agent/monitor/config"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
	"github.com/volcengine/volc-sdk-golang/base"
)

const (
	// DefaultRegion default region of instance
	DefaultRegion = "x"
	// ServiceVersion20180801 .
	ServiceVersion20180801 = "x"
	// ServiceName .
	ServiceName = utils.CloudMonitorService
)

var (
	// ServiceInfo .
	ServiceInfo = &base.ServiceInfo{
		Timeout: 7 * time.Second,
		Host:    "x",
		Header: http.Header{
			"Accept":           []string{"application/json"},
			"X-Tsdb-Namespace": []string{"ecs"},
			"X-Tsdb-Region":    []string{config.GetMonitorConfig().Region},
		},
	}
	// APIInfoList .
	APIInfoList = map[string]*base.ApiInfo{
		"SendMetricData": {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{"Write"},
				"Version": []string{ServiceVersion20180801},
			},
		},
	}
)

// DefaultInstance .
var DefaultInstance = NewInstance()

// Metric IAM .
type Metric struct {
	Client *base.Client
}

// NewInstance 创建一个实例
func NewInstance() *Metric {
	instance := &Metric{}
	instance.Client = base.NewClient(ServiceInfo, APIInfoList)
	instance.Client.ServiceInfo.Credentials.Service = ServiceName
	region := config.GetMonitorConfig().Region
	if region != "" {
		instance.Client.ServiceInfo.Credentials.Region = region
	} else {
		instance.Client.ServiceInfo.Credentials.Region = DefaultRegion
	}

	return instance
}

// GetServiceInfo interface
func (m *Metric) GetServiceInfo() *base.ServiceInfo {
	return ServiceInfo
}

// GetAPIInfo interface
func (m *Metric) GetAPIInfo(api string) *base.ApiInfo {
	if apiInfo, ok := APIInfoList[api]; ok {
		return apiInfo
	}
	return nil
}

// SetRegion SetHost .
func (m *Metric) SetRegion(region string) {
	ServiceInfo.Credentials.Region = region
}

// SetHost .
func (m *Metric) SetHost(host string) {
	m.Client.ServiceInfo.Host = host
}

// SetSchema .
func (m *Metric) SetSchema(schema string) {
	m.Client.ServiceInfo.Scheme = schema
}
