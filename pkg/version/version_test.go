package version

import (
	"testing"

	. "github.com/moul/advanced-ssh-config/vendor/github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("Testing version", t, func() {
		So(VERSION, ShouldNotEqual, "")
		So(GITCOMMIT, ShouldNotEqual, "")
	})
}
