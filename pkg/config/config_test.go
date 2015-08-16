package config

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func dummyConfig() *Config {
	config := New()
	config.Hosts["toto"] = Host{
		Host: "1.2.3.4",
	}
	config.Hosts["titi"] = Host{
		Host: "tata",
		Port: 23,
		User: "moul",
	}
	config.Defaults = Host{
		Port: 22,
		User: "root",
	}
	return config
}

func TestNew(t *testing.T) {
	Convey("Testing New()", t, func() {
		config := New()

		So(len(config.Hosts), ShouldEqual, 0)
		So(config.Defaults.Port, ShouldEqual, uint(0))
		So(config.Defaults.Host, ShouldEqual, "")
		So(config.Defaults.User, ShouldEqual, "")
	})
}

func TestConfig(t *testing.T) {
	Convey("Testing Config", t, func() {
		config := dummyConfig()

		So(len(config.Hosts), ShouldEqual, 2)
		So(config.Hosts["toto"].Host, ShouldEqual, "1.2.3.4")
		So(config.Defaults.Port, ShouldEqual, uint(22))
	})
}

func TestConfig_LoadConfig(t *testing.T) {
	Convey("Testing Config.LoadConfig", t, func() {

		config := New()
		err := config.LoadConfig(strings.NewReader(`
hosts:
  aaa:
    host: 1.2.3.4
  bbb:
    port: 21
  ccc:
    host: 5.6.7.8
    port: 24
    user: toor
defaults:
  port: 22
  user: root
`))
		So(err, ShouldBeNil)
		So(len(config.Hosts), ShouldEqual, 3)
		So(config.Hosts["aaa"].Host, ShouldEqual, "1.2.3.4")
		So(config.Hosts["aaa"].Port, ShouldEqual, uint(0))
		So(config.Hosts["aaa"].User, ShouldEqual, "")
		So(config.Hosts["bbb"].Host, ShouldEqual, "")
		So(config.Hosts["bbb"].Port, ShouldEqual, uint(21))
		So(config.Hosts["bbb"].User, ShouldEqual, "")
		So(config.Hosts["ccc"].Host, ShouldEqual, "5.6.7.8")
		So(config.Hosts["ccc"].Port, ShouldEqual, uint(24))
		So(config.Hosts["ccc"].User, ShouldEqual, "toor")
		So(config.Defaults.Port, ShouldEqual, uint(22))
		So(config.Defaults.User, ShouldEqual, "root")
	})
}
