package collectors

import (
	"strings"

	"github.com/shirou/gopsutil/net"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/config"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
	"go.uber.org/zap"
)

// NicEthPrefixs store some other style of network interface name.
var NicEthPrefixs = []string{"eth"}

// NetStates is the type for store net state.
type NetStates struct {
	byteSent        float64
	byteRecv        float64
	packetSent      float64
	packetRecv      float64
	errIn           float64
	errOut          float64
	dropIn          float64
	dropOut         float64
	collectTime     int64
	uptimeInSeconds int64
}

// NetCollector is the collector type for net metric.
type NetCollector struct {
	LastStates *NetStates
}

// Collect implement the net Collector.
func (n *NetCollector) Collect(collectTime int64) *model.InputMetric {
	var NetworkInErrorPackages, NetworkOutErrorPackages, NetworkInDrop, NetworkOutDrop float64
	deltaTime := float64(config.DefaultMetricDeltaDataTimeInSecond)
	netStates, err := net.IOCounters(true)

	if nil != err {
		logs.GetLogger().Error("get net io count error", zap.Error(err))
		return nil
	}

	allStats := getIOCountersAll(netStates)
	currStates := &NetStates{
		byteSent:   float64(allStats.BytesSent),
		byteRecv:   float64(allStats.BytesRecv),
		packetSent: float64(allStats.PacketsSent),
		packetRecv: float64(allStats.PacketsRecv),
		errIn:      float64(allStats.Errin),
		errOut:     float64(allStats.Errout),
		dropIn:     float64(allStats.Dropin),
		dropOut:    float64(allStats.Dropout),

		collectTime: collectTime,
	}
	currStates.uptimeInSeconds, err = utils.GetUptimeInSeconds()
	if err != nil {
		logs.GetLogger().Error("exec GetUptimeInSeconds error", zap.Error(err))
	}

	if n.LastStates == nil {
		n.LastStates = currStates
		return nil
	}

	deltaTimeUsingCT := float64(currStates.collectTime-n.LastStates.collectTime) / 1000
	if currStates.uptimeInSeconds != -1 && n.LastStates.uptimeInSeconds != -1 {
		deltaTime = float64(currStates.uptimeInSeconds - n.LastStates.uptimeInSeconds)
	} else if deltaTimeUsingCT > 0 {
		deltaTime = deltaTimeUsingCT
	}

	totalSentPacket := utils.Float64From32Bits(currStates.packetSent - n.LastStates.packetSent)
	totalRecvPacket := utils.Float64From32Bits(currStates.packetRecv - n.LastStates.packetRecv)
	NetworkInRate := utils.Float64From32Bits(currStates.byteRecv-n.LastStates.byteRecv) / deltaTime * model.ByteToBit
	NetworkOutRate := utils.Float64From32Bits(currStates.byteSent-n.LastStates.byteSent) / deltaTime * model.ByteToBit
	NetworkOutPackages := totalSentPacket / deltaTime
	NetworkInPackages := totalRecvPacket / deltaTime

	// there's hard coding but I don't think it's so hard.
	conns, _ := net.Connections("tcp")
	NetTCPConnection := float64(len(conns))

	if totalRecvPacket != 0 {
		NetworkInErrorPackages = model.ToPercent * utils.Float64From32Bits(currStates.errIn-n.LastStates.errIn) /
			totalRecvPacket / deltaTime
		NetworkInDrop = model.ToPercent * utils.Float64From32Bits(currStates.dropIn-n.LastStates.dropIn) /
			totalRecvPacket / deltaTime
	}
	if totalSentPacket != 0 {
		NetworkOutErrorPackages = model.ToPercent * utils.Float64From32Bits(currStates.errOut-n.LastStates.errOut) /
			totalSentPacket / deltaTime
		NetworkOutDrop = model.ToPercent * utils.Float64From32Bits(currStates.dropOut-n.LastStates.dropOut) /
			totalSentPacket / deltaTime
	}

	n.LastStates = currStates

	metricsDatas := []model.Metric{
		{
			MetricName:  "NetworkOutRate",
			MetricValue: NetworkOutRate,
		},
		{
			MetricName:  "NetworkInRate",
			MetricValue: NetworkInRate,
		},
		{
			MetricName:  "NetworkOutPackages",
			MetricValue: NetworkOutPackages,
		},
		{
			MetricName:  "NetworkInPackages",
			MetricValue: NetworkInPackages,
		},
		{
			MetricName:  "NetworkInErrorPackages",
			MetricValue: NetworkInErrorPackages,
		},
		{
			MetricName:  "NetworkOutErrorPackages",
			MetricValue: NetworkOutErrorPackages,
		},
		{
			MetricName:  "NetworkInDrop",
			MetricValue: NetworkInDrop,
		},
		{
			MetricName:  "NetworkOutDrop",
			MetricValue: NetworkOutDrop,
		},
		{
			MetricName:  "NetTcpConnection",
			MetricValue: NetTCPConnection,
		},
	}

	return &model.InputMetric{
		Data:        metricsDatas,
		Type:        "network",
		CollectTime: collectTime,
	}
}

func getIOCountersAll(n []net.IOCountersStat) net.IOCountersStat {
	r := net.IOCountersStat{
		Name: "all",
	}
	for _, nic := range n {
		if !isNeedToCal(nic.Name) {
			continue
		}
		r.BytesRecv += nic.BytesRecv
		r.PacketsRecv += nic.PacketsRecv
		r.Errin += nic.Errin
		r.Dropin += nic.Dropin
		r.BytesSent += nic.BytesSent
		r.PacketsSent += nic.PacketsSent
		r.Errout += nic.Errout
		r.Dropout += nic.Dropout
	}

	return r
}

func isNeedToCal(nic string) bool {
	for _, nicEthPrefix := range NicEthPrefixs {
		if strings.HasPrefix(nic, nicEthPrefix) {
			return true
		}
	}
	return false

}
