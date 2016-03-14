package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// BoolVal returns a boolean matching a configuration string
func TestBoolVal(t *testing.T) {
	Convey("Testing BoolVal", t, func() {
		trueValues := []string{"yes", "ok", "true", "1", "enabled", "True", "TRUE", "YES", "Yes"}
		falseValues := []string{"no", "0", "false", "False", "FALSE", "disabled"}
		for _, val := range trueValues {
			So(BoolVal(val), ShouldBeTrue)
		}
		for _, val := range falseValues {
			So(BoolVal(val), ShouldBeFalse)
		}
	})
}
