package utils

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHomeDir(t *testing.T) {
	Convey("Testing GetHomeDir", t, func() {
		oldHome := os.Getenv("HOME")
		oldUserProfile := os.Getenv("USERPROFILE")

		So(os.Setenv("HOME", "/a/b/c"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		So(GetHomeDir(), ShouldEqual, "/a/b/c")

		So(os.Setenv("HOME", "/a/b/d"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		So(GetHomeDir(), ShouldEqual, "/a/b/d")

		So(os.Setenv("HOME", "/a/b/d"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/e"), ShouldBeNil)
		So(GetHomeDir(), ShouldEqual, "/a/b/d")

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/f"), ShouldBeNil)
		So(GetHomeDir(), ShouldEqual, "/a/b/f")

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/g"), ShouldBeNil)
		So(GetHomeDir(), ShouldEqual, "/a/b/g")

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		So(GetHomeDir(), ShouldEqual, "")

		So(os.Setenv("HOME", oldHome), ShouldBeNil)
		So(os.Setenv("USERPROFILE", oldUserProfile), ShouldBeNil)
	})
}

func TestExpandUser(t *testing.T) {
	Convey("Testing ExpandUser", t, func() {
		oldHome := os.Getenv("HOME")
		oldUserProfile := os.Getenv("USERPROFILE")

		So(os.Setenv("HOME", "/a/b/c"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err := ExpandUser("~/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/a/b/d"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, "/a/b/d/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/a/b/d"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/e"), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, "/a/b/d/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/f"), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, "/a/b/f/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/g"), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, "/a/b/g/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, "")
		So(err, ShouldNotBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/e/f"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/e/g"), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/e/h"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/e/i"), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", oldHome), ShouldBeNil)
		So(os.Setenv("USERPROFILE", oldUserProfile), ShouldBeNil)
	})
}
