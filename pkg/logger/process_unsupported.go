// +build openbsd netbsd freebsd

package logger

import "go.uber.org/zap/zapcore"

// LogLevelFromParentSSHProcess inspects parent `ssh` process for eventual passed `-v` flags.
func LogLevelFromParentSSHProcess() (zapcore.Level, error) {
	return zapcore.WarnLevel, nil
}
