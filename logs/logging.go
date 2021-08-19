package logs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	error2 "github.com/volcengine/cloud-monitor-agent/error"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger              *zap.Logger
	logLevel            string
	logDir              string
	singleFileMaxSizeMB int
	maxBackups          int
	logsKeepDay         int
	loggerSync          sync.Once
)

// Config define the config of log
type Config struct {
	LogLevel            string
	LogDir              string
	SingleFileMaxSizeMB int
	MaxBackups          int
	LogsKeepDay         int
}

// initLogDir init location of log.
func initLogDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		}
		return err
	}

	if info.IsDir() {
		return nil
	}

	return error2.Errors.DirPathError
}

func initLogger() {
	pwd := GetCurrentDirectory()
	file, err := os.Open(pwd + "/conf_logs.json")
	if err != nil {
		GetLogger().Error("Open conf_logs.json configuration file", zap.Error(err))
		return
	}

	decoder := json.NewDecoder(file)
	logsConfig := Config{}
	err = decoder.Decode(&logsConfig)
	if err != nil {
		GetLogger().Error("Parse conf_logs.json failed", zap.Error(err))
		return
	}

	logLevel = logsConfig.LogLevel
	logDir = logsConfig.LogDir
	singleFileMaxSizeMB = logsConfig.SingleFileMaxSizeMB
	maxBackups = logsConfig.MaxBackups
	logsKeepDay = logsConfig.LogsKeepDay

	var level zapcore.Level
	var encoder zapcore.Encoder
	var appName string

	if !filepath.IsAbs(logDir) {
		logDir, err = filepath.Abs(logDir)
		if err != nil {
			panic(err)
		}
	}
	executableName, err := os.Executable()
	if err != nil {
		appName = "unknown"
	} else {
		appName = filepath.Base(executableName)
	}
	initLogDir(logDir)

	level.Set(logLevel)
	levelFunc := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		if l >= level {
			return true
		}
		return false
	})

	encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logDir, appName+".log"),
		MaxSize:    singleFileMaxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     logsKeepDay,
	})

	c := zapcore.NewCore(
		encoder,
		w,
		levelFunc,
	)

	logger = zap.New(c)
	if err != nil {
		panic(err)
	}

}

// GetLogger return instance of logger
func GetLogger() *zap.Logger {
	loggerSync.Do(initLogger)
	return logger
}

// GetCurrentDirectory return dir of project
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		GetLogger().Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}
