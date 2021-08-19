package config

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	error2 "github.com/volcengine/cloud-monitor-agent/error"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/monitor/utils"
	"go.uber.org/zap"
)

// MonitorConfig recode all of the config about agent.
type MonitorConfig struct {
	Endpoint    string
	MetaService string
	Region      string
}

// PluginConfig recode plugin status
type PluginConfig struct {
	GPUCollectorOpen bool
}

var (
	monitorConfig *MonitorConfig
	pluginConfig  *PluginConfig
	// DefaultMetricDeltaDataTimeInSecond not allowed the user modify
	DefaultMetricDeltaDataTimeInSecond = utils.CronJobTimeSecond
)

// loadConfig load config from file.
func loadConfig(confName string, conf interface{}) (interface{}, error) {
	pwd := logs.GetCurrentDirectory()
	file, err := os.Open(pwd + "/" + confName)
	if err != nil {
		logs.GetLogger().Error("Open monitor configuration file error", zap.Error(err))
		return nil, error2.Errors.NoConfigFileFound
	}
	defer file.Close()

	decoder := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		logs.GetLogger().Error("Parsing monitor configuration file error", zap.Error(err))
		return nil, error2.Errors.ConfigFileValidationError
	}
	logs.GetLogger().Info("Successfully loaded monitor configuration file")
	return conf, nil
}

// GetMonitorConfig get config.
func GetMonitorConfig() *MonitorConfig {
	return monitorConfig
}

// GetPluginConfig get config.
func GetPluginConfig() *PluginConfig {
	return pluginConfig
}

func init() {
	var ok bool
	monitorConf, err1 := loadConfig(utils.ConfMonitorName, monitorConfig)
	pluginConf, err2 := loadConfig(utils.ConfPluginName, pluginConfig)
	if err1 != nil || err2 != nil {
		logs.GetLogger().Error("open config file error")
		panic(error2.Errors.OpenConfigFileError)
	}

	if monitorConfig, ok = monitorConf.(*MonitorConfig); !ok {
		logs.GetLogger().Error("type conversion error")
		panic(error2.Errors.CastTypeError)
	}

	if pluginConfig, ok = pluginConf.(*PluginConfig); !ok {
		logs.GetLogger().Error("type conversion error")
		panic(error2.Errors.CastTypeError)
	}
}
