package logger

import "github.com/Sirupsen/logrus"

var Logger = logrus.New()

// SetLevel sets the logging level
func SetLevel(level logrus.Level) {
	// Logger.mu.Lock()
	// defer Logger.mu.Unlock()
	Logger.Level = level
}

// Options allows to customize logger behavior
type Options struct {
	Level         logrus.Level
	InspectParent bool
}

// SetupLogging configures the logger based on user input and parent process configuration (looks for `ssh -v`)
func SetupLogging(options Options) {
	level := options.Level

	if options.InspectParent {
		parentLevel, err := GetLoggingLevelByInspectingParent()
		if err != nil {
			Logger.Debugf("Failed to inspect parent process: %v", err)
		} else if parentLevel > level {
			level = parentLevel
		}
	}

	SetLevel(level)
}
