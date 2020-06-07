package version

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("Testing version", t, func() {
		So(Version, ShouldNotEqual, "")
		So(VcsRef, ShouldNotEqual, "")
	})
}
