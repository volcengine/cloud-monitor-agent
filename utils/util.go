package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"go.uber.org/zap"
)

var (
	workingPath string
)

// GetCurrTSInNano return the time in Nano.
func GetCurrTSInNano() int64 {
	return time.Now().UnixNano()
}

// GetWorkingPath get cloud-monitor-agent woring path.
func GetWorkingPath() string {
	var err error

	if workingPath == "" {
		workingPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logs.GetLogger().Error("Get working path path failed", zap.Error(err))
			return ""
		}
	}

	return workingPath
}

// GetMainThreadPidFilePath .
func GetMainThreadPidFilePath() string {
	return filepath.Join(GetWorkingPath(), AgentPIDFileName)
}

// GetMainProcessPid will return the pid of main process.
func GetMainProcessPid() (int, error) {
	content, err := ioutil.ReadFile(GetMainThreadPidFilePath())
	// ignore the error if file is not exist.
	if err != nil {
		return 0, errors.Wrap(err, "open main thread pid file failed")
	}

	num, err := strconv.Atoi(strings.Split(string(content), "\n")[0])
	if err != nil {
		return 0, errors.Wrap(err, "content format of PID file is error")
	}

	return num, nil
}
