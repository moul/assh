// +build !openbsd
// +build !freebsd
// +build !netbsd

package logger

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/process"
)

func GetLoggingLevelByInspectingParent() (logrus.Level, error) {
	ppid := os.Getppid()
	process, err := process.NewProcess(int32(ppid))
	if err != nil {
		return logrus.WarnLevel, err
	}

	cmdline, err := process.Cmdline()
	if err != nil {
		return logrus.WarnLevel, err
	}

	if strings.Contains(cmdline, "-vv") {
		return logrus.DebugLevel, nil
	} else if strings.Contains(cmdline, "-v") {
		return logrus.InfoLevel, nil
	} else if strings.Contains(cmdline, "-q") {
		return logrus.ErrorLevel, nil
	}
	return logrus.WarnLevel, nil
}
