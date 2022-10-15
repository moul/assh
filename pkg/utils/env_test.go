package utils

import (
	"os"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEscapeSpaces(t *testing.T) {
	Convey("Testing EscapeSpaces", t, func() {
		So(EscapeSpaces("foo bar"), ShouldEqual, "foo\\ bar")
		So(EscapeSpaces("/a/b c/d"), ShouldEqual, "/a/b\\ c/d")
	})
}

func TestExpandEnvSafe(t *testing.T) {
	Convey("Testing ExpandEnvSafe", t, func() {
		So(os.Setenv("FOO", "bar"), ShouldBeNil)
		So(ExpandEnvSafe("/a/$FOO/c"), ShouldEqual, "/a/bar/c")
		So(ExpandEnvSafe("/a/${FOO}/c"), ShouldEqual, "/a/bar/c")
		So(ExpandEnvSafe("/a/${FOO}/c/$FOO"), ShouldEqual, "/a/bar/c/bar")
		So(ExpandEnvSafe("/a/$(FOO)/c"), ShouldEqual, "/a/$(FOO)/c")
		So(ExpandEnvSafe(""), ShouldEqual, "")
	})
}

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
	expected := "/a/b/c/test"
	expectedEscaped := "/a/b/c/test\\ dir"

	if runtime.GOOS == "windows" {
		expected = "\\a\\b\\c\\test"
		expectedEscaped = "\\a\\b\\c\\test\\ dir"
	}

	Convey("Testing ExpandUser", t, func() {
		oldHome := os.Getenv("HOME")
		oldUserProfile := os.Getenv("USERPROFILE")

		So(os.Setenv("HOME", "/a/b/c"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err := ExpandUser("~/test dir")
		So(dir, ShouldEqual, expectedEscaped)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/a/b/c"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/a/b/c"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/a/b/c"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/e"), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/a/b/c"), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("~/test")
		So(dir, ShouldEqual, "")
		So(err, ShouldNotBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/e/f"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", ""), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", ""), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/e/g"), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", "/e/h"), ShouldBeNil)
		So(os.Setenv("USERPROFILE", "/e/i"), ShouldBeNil)
		dir, err = ExpandUser("/a/b/c/test")
		So(dir, ShouldEqual, expected)
		So(err, ShouldBeNil)

		So(os.Setenv("HOME", oldHome), ShouldBeNil)
		So(os.Setenv("USERPROFILE", oldUserProfile), ShouldBeNil)
	})
}
