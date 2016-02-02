package logger

import (
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("Testing Logger", t, func() {
		So(Logger, ShouldNotBeNil)
	})
}

func TestSetupLogging(t *testing.T) {
	Convey("Testing SetupLogging()", t, func() {
		Reset(func() {
			Logger.Level = logrus.InfoLevel
		})
		Convey("InspectParent=false", func() {
			So(Logger.Level, ShouldEqual, logrus.InfoLevel)
			options := LoggerOptions{
				Level:         logrus.WarnLevel,
				InspectParent: false,
			}
			SetupLogging(options)
			So(Logger.Level, ShouldEqual, logrus.WarnLevel)
		})
		Convey("InspectParent=true", func() {
			So(Logger.Level, ShouldEqual, logrus.InfoLevel)
			options := LoggerOptions{
				Level:         logrus.WarnLevel,
				InspectParent: true,
			}
			SetupLogging(options)
			// FIXME: mock process
			// So(Logger.Level, ShouldEqual, logrus.InfoLevel)
		})
	})
}

func TestGetLoggingLevelByInspectingParent(t *testing.T) {
	Convey("Testing GetLoggingLevelByInspectingParent()", t, func() {
		_, err := GetLoggingLevelByInspectingParent()
		So(err, ShouldBeNil)
		// FIXME: mock process
		// So(level, ShouldEqual, logrus.InfoLevel)
	})
}
