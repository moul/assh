// +build !openbsd
// +build !freebsd
// +build !netbsd

package logger

import (
	"os"
	"strings"

	"github.com/shirou/gopsutil/process"
	"go.uber.org/zap/zapcore"
)

// LogLevelFromParentSSHProcess inspects parent `ssh` process for eventual passed `-v` flags.
func LogLevelFromParentSSHProcess() (zapcore.Level, error) {
	// FIXME: check if parent process is `ssh`
	ppid := os.Getppid()
	process, err := process.NewProcess(int32(ppid))
	if err != nil {
		return zapcore.WarnLevel, err
	}

	cmdline, err := process.Cmdline()
	if err != nil {
		return zapcore.WarnLevel, err
	}

	if strings.Contains(cmdline, "-vv") {
		return zapcore.DebugLevel, nil
	} else if strings.Contains(cmdline, "-v") {
		return zapcore.InfoLevel, nil
	} else if strings.Contains(cmdline, "-q") {
		return zapcore.ErrorLevel, nil
	}
	return zapcore.WarnLevel, nil
}
