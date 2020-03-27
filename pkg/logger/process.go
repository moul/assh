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

	switch {
	case strings.Contains(cmdline, "-vv"):
		return zapcore.DebugLevel, nil
	case strings.Contains(cmdline, "-v"):
		return zapcore.InfoLevel, nil
	case strings.Contains(cmdline, "-q"):
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.WarnLevel, nil
	}
}
