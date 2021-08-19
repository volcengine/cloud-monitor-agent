package manager

import (
	"github.com/nightlyone/lockfile"
	"github.com/pkg/errors"
	"github.com/volcengine/cloud-monitor-agent/logs"
	"github.com/volcengine/cloud-monitor-agent/utils"
	"go.uber.org/zap"
)

// SingleInstanceCheck check whether the agent process is running before start
// error indicates there is another agent process is running or sth.
func SingleInstanceCheck() (err error) {
	// is there need to move cloudMonitorAgent.pid to /run/cloudMonitorAgent.pid?
	lock, err := lockfile.New(utils.GetMainThreadPidFilePath())
	if err != nil {
		return errors.Wrap(err, "Cannot init pid file")
	}

	err = lock.TryLock()
	if err != nil {
		p, _ := lock.GetOwner()
		logs.GetLogger().Debug("GetOwner successfully and pid in file ", zap.Any("pid", p.Pid))

		return errors.Wrap(err, "Cannot lock")
	}
	return nil
}
