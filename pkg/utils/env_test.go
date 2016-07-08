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

		os.Setenv("HOME", "/a/b/c")
		os.Setenv("USERPROFILE", "")
		So(GetHomeDir(), ShouldEqual, "/a/b/c")

		os.Setenv("HOME", "/a/b/d")
		os.Setenv("USERPROFILE", "")
		So(GetHomeDir(), ShouldEqual, "/a/b/d")

		os.Setenv("HOME", "/a/b/d")
		os.Setenv("USERPROFILE", "/a/b/e")
		So(GetHomeDir(), ShouldEqual, "/a/b/d")

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "/a/b/f")
		So(GetHomeDir(), ShouldEqual, "/a/b/f")

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "/a/b/g")
		So(GetHomeDir(), ShouldEqual, "/a/b/g")

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "")
		So(GetHomeDir(), ShouldEqual, "")

		os.Setenv("HOME", oldHome)
		os.Setenv("USERPROFILE", oldUserProfile)
	})
}

func TestExpandUser(t *testing.T) {
	Convey("Testing expandUser", t, func() {
		oldHome := os.Getenv("HOME")
		oldUserProfile := os.Getenv("USERPROFILE")

		os.Setenv("HOME", "/a/b/c")
		os.Setenv("USERPROFILE", "")
		dir, err := expandUser("~/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "/a/b/d")
		os.Setenv("USERPROFILE", "")
		dir, err = expandUser("~/test")
		So(dir, ShouldEqual, "/a/b/d/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "/a/b/d")
		os.Setenv("USERPROFILE", "/a/b/e")
		dir, err = expandUser("~/test")
		So(dir, ShouldEqual, "/a/b/d/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "/a/b/f")
		dir, err = expandUser("~/test")
		So(dir, ShouldEqual, "/a/b/f/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "/a/b/g")
		dir, err = expandUser("~/test")
		So(dir, ShouldEqual, "/a/b/g/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "")
		dir, err = expandUser("~/test")
		So(dir, ShouldEqual, "")
		So(err, ShouldNotBeNil)

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "")
		dir, err = expandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "/e/f")
		os.Setenv("USERPROFILE", "")
		dir, err = expandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "/e/g")
		dir, err = expandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", "/e/h")
		os.Setenv("USERPROFILE", "/e/i")
		dir, err = expandUser("/a/b/c/test")
		So(dir, ShouldEqual, "/a/b/c/test")
		So(err, ShouldBeNil)

		os.Setenv("HOME", oldHome)
		os.Setenv("USERPROFILE", oldUserProfile)
	})
}
