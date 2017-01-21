package logger

import "github.com/Sirupsen/logrus"

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
