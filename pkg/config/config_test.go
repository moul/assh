package config

import (
	"strings"
	"testing"

	. "github.com/moul/advanced-ssh-config/vendor/github.com/smartystreets/goconvey/convey"
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
	config.Hosts["*.ddd"] = Host{
		Host: "1.3.5.7",
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

		So(len(config.Hosts), ShouldEqual, 3)
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
  "*.ddd":
    host: 1.3.5.7
defaults:
  port: 22
  user: root
`))
		So(err, ShouldBeNil)
		So(len(config.Hosts), ShouldEqual, 4)
		So(config.Hosts["aaa"].Host, ShouldEqual, "1.2.3.4")
		So(config.Hosts["aaa"].Port, ShouldEqual, uint(0))
		So(config.Hosts["aaa"].User, ShouldEqual, "")
		So(config.Hosts["bbb"].Host, ShouldEqual, "")
		So(config.Hosts["bbb"].Port, ShouldEqual, uint(21))
		So(config.Hosts["bbb"].User, ShouldEqual, "")
		So(config.Hosts["ccc"].Host, ShouldEqual, "5.6.7.8")
		So(config.Hosts["ccc"].Port, ShouldEqual, uint(24))
		So(config.Hosts["ccc"].User, ShouldEqual, "toor")
		So(config.Hosts["*.ddd"].Host, ShouldEqual, "1.3.5.7")
		So(config.Hosts["*.ddd"].Port, ShouldEqual, uint(0))
		So(config.Hosts["*.ddd"].User, ShouldEqual, "")
		So(config.Defaults.Port, ShouldEqual, uint(22))
		So(config.Defaults.User, ShouldEqual, "root")
	})
}

func TestConfig_JsonSring(t *testing.T) {
	Convey("Testing Config.JsonString", t, func() {
		config := dummyConfig()
		expected := `{
  "hosts": {
    "*.ddd": {
      "Name": "",
      "host": "1.3.5.7"
    },
    "titi": {
      "Name": "",
      "host": "tata",
      "user": "moul",
      "port": 23
    },
    "toto": {
      "Name": "",
      "host": "1.2.3.4"
    }
  },
  "defaults": {
    "Name": "",
    "user": "root",
    "port": 22
  },
  "includes": null
}`
		json, err := config.JsonString()
		So(err, ShouldBeNil)
		So(string(json), ShouldEqual, expected)
	})
}

func TestConfig_getHostByName(t *testing.T) {
	Convey("Testing Config.getHostByName", t, func() {

		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.getHostByName("titi", false)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")

			host, err = config.getHostByName("titi", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")

			host, err = config.getHostByName("dontexists", false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("dontexists", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "dontexists")

			host, err = config.getHostByName("regex.ddd", false)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")

			host, err = config.getHostByName("regex.ddd", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
		})

		Convey("With gateway", func() {
			host, err = config.getHostByName("titi/gateway", false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("titi/gateway", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByName("dontexists/gateway", false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("dontexists/gateway", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "dontexists/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByName("regex.ddd/gateway", false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("regex.ddd/gateway", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd/gateway")
			So(host.Host, ShouldNotEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})
	})
}

func TestConfig_getHostByPath(t *testing.T) {
	Convey("Testing Config.getHostByPath", t, func() {

		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.getHostByPath("titi", false)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("titi", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("dontexists", false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByPath("dontexists", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("regex.ddd", false)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("regex.ddd", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host, err = config.getHostByPath("titi/gateway", false)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("titi/gateway", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("dontexists/gateway", false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByPath("dontexists/gateway", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("regex.ddd/gateway", false)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("regex.ddd/gateway", true)
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})
	})
}

func TestConfig_GetHost(t *testing.T) {
	Convey("Testing Config.GetHost", t, func() {

		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.GetHost("titi")
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.GetHost("dontexists")
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.GetHost("regex.ddd")
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host, err = config.GetHost("titi/gateway")
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.GetHost("dontexists/gateway")
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			// FIXME: check if this is a normal behavior
			host, err = config.GetHost("regex.ddd/gateway")
			So(err, ShouldBeNil)
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})
	})
}

func TestConfig_GetHostSafe(t *testing.T) {
	Convey("Testing Config.GetHostSafe", t, func() {

		config := dummyConfig()
		var host *Host

		Convey("Without gateway", func() {
			host = config.GetHostSafe("titi")
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetHostSafe("dontexists")
			So(host.Name, ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetHostSafe("regex.ddd")
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host = config.GetHostSafe("titi/gateway")
			So(host.Name, ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host = config.GetHostSafe("dontexists/gateway")
			So(host.Name, ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 1)

			host = config.GetHostSafe("regex.ddd/gateway")
			So(host.Name, ShouldEqual, "regex.ddd")
			So(host.Host, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})
	})
}
