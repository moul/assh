// +build openbsd netbsd freebsd

package logger

import "github.com/sirupsen/logrus"

func GetLoggingLevelByInspectingParent() (logrus.Level, error) {
	return logrus.WarnLevel, nil
}
