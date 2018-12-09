package logger

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// FIXME: test MustLogLevel

func TestLogLevelFromParentSSHProcess(t *testing.T) {
	Convey("Testing LogLevelFromParentSSHProcess()", t, func() {
		_, err := LogLevelFromParentSSHProcess()
		So(err, ShouldBeNil)
		// FIXME: mock process
		// So(level, ShouldEqual, zapcore.InfoLevel)
	})
}
