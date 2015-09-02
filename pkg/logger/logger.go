package logger

import "github.com/moul/advanced-ssh-config/vendor/github.com/Sirupsen/logrus"

var Logger = logrus.New()

func LoggerSetLevel(level logrus.Level) {
	// Logger.mu.Lock()
	// defer Logger.mu.Unlock()
	Logger.Level = level
}
