package logger

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/process"
)

var Logger = logrus.New()

func LoggerSetLevel(level logrus.Level) {
	// Logger.mu.Lock()
	// defer Logger.mu.Unlock()
	Logger.Level = level
}

type LoggerOptions struct {
	Level         logrus.Level
	InspectParent bool
}

func SetupLogging(options LoggerOptions) {
	level := options.Level

	if options.InspectParent {
		parentLevel, err := GetLoggingLevelByInspectingParent()
		if err != nil {
			Logger.Debugf("Failed to inspect parent process: %v", err)
		} else if parentLevel > level {
			level = parentLevel
		}
	}

	LoggerSetLevel(level)
}

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
	}
	return logrus.WarnLevel, nil
}
