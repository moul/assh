package logger // import "moul.io/assh/pkg/logger"

import "go.uber.org/zap/zapcore"

// MustLogLevel returns a log level based on both user input and parent SSH process
func MustLogLevel(debug, verbose bool) zapcore.Level {
	parentLevel, err := LogLevelFromParentSSHProcess()
	if err != nil {
		parentLevel = zapcore.WarnLevel
	}
	asshLevel := zapcore.WarnLevel
	switch {
	case debug:
		asshLevel = zapcore.DebugLevel
	case verbose:
		asshLevel = zapcore.InfoLevel
	}
	if parentLevel > asshLevel {
		return parentLevel
	}
	return asshLevel
}
