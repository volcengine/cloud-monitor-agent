package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	errors2 "github.com/volcengine/cloud-monitor-agent/error"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/config"
	"github.com/volcengine/cloud-monitor-agent/monitor/model"
	"github.com/volcengine/cloud-monitor-agent/monitor/report/sdk/cloudmonitor"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
	"go.uber.org/zap"
)

// MetaData is info of gataway.
type MetaData struct {
	ExpiredTime     string
	CurrentTime     string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// getToken is return sts of ECS or BES.
func getToken() (MetaData, error) {
	var metaData MetaData
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := http.Client{Timeout: 1200 * time.Millisecond, Transport: &transport}
	resp, err := client.Get(config.GetMonitorConfig().MetaService)

	if err != nil {
		return metaData, errors.Wrap(err, "metaService unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("metaService StatusCode error, code: %d", resp.StatusCode)

		return metaData, errors.Wrap(errors2.Errors.MetaServiceError, errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &metaData); err != nil {
		return metaData, errors.Wrap(err, "unmarshal json error")
	}

	return metaData, nil
}

// flushToken set sts to sdk client.
func flushToken(url string) error {
	metaData, err := getToken()
	if err != nil {
		return err
	}

	cloudmonitor.DefaultInstance.Client.SetAccessKey(metaData.AccessKeyID)
	cloudmonitor.DefaultInstance.Client.SetSecretKey(metaData.SecretAccessKey)
	cloudmonitor.DefaultInstance.Client.SetSessionToken(metaData.SessionToken)
	cloudmonitor.DefaultInstance.SetHost(url)

	return nil
}

// wrapMetricDataToReq switch MetricsDataArr to the format required for final sending.
func wrapMetricDataToReq(datas model.MetricDataArr) *cloudmonitor.SendDataRequest {
	var metricDatas string

	for _, data := range datas {
		var s string
		var dmsStr string
		for _, dimension := range data.Metric.Dimensions {
			dmsStr += dimension.Name + "=" + dimension.Value + ","
		}
		dmsStr = dmsStr[:len(dmsStr)-1]
		s = utils.MeasurementPrefix + "_" + data.Metric.MetricName + "," + dmsStr + " value=" +
			strconv.FormatFloat(data.Value, 'f', -1, 64) + " " +
			strconv.FormatInt(int64(data.CollectTime), 10) + "\n"
		metricDatas = metricDatas + s
	}

	return &cloudmonitor.SendDataRequest{DataLines: metricDatas}
}

// SendMetricData send metric data to the endpoint.
func SendMetricData(url string, data *model.InputMetric) {
	defer func() {
		if err := recover(); err != nil {
			logs.GetLogger().Panic("send MetricData task panic", zap.String("reason", fmt.Sprintf("%v", err)))
		}
	}()

	if err := flushToken(url); err != nil {
		logs.GetLogger().Error("flushToken failed", zap.Error(err))
		return
	}

	metricDataArr := model.BuildMetricData(data)
	resp, retCode, err := cloudmonitor.DefaultInstance.Send(wrapMetricDataToReq(metricDataArr))

	if err != nil && retCode != http.StatusNoContent && retCode != http.StatusOK {
		logs.GetLogger().Error("SendMetricData fail, the error", zap.Error(err))
	}
	logs.GetLogger().Debug("resp from Volc_InfluxDB_Proxy", zap.Any("resp", resp))

}
