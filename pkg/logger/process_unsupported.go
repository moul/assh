// +build openbsd

package logger

import "github.com/Sirupsen/logrus"

func GetLoggingLevelByInspectingParent() (logrus.Level, error) {
	return logrus.WarnLevel, nil
}
