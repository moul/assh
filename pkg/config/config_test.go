package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/moul/advanced-ssh-config/vendor/github.com/smartystreets/goconvey/convey"
)

var (
	configExample string = `
hosts:
  aaa:
    HostName: 1.2.3.4
  bbb:
    Port: 21
  ccc:
    HostName: 5.6.7.8
    Port: 24
    User: toor
  "*.ddd":
    HostName: 1.3.5.7
defaults:
  Port: 22
  User: root
`
)

func TestNew(t *testing.T) {
	Convey("Testing New()", t, func() {
		config := New()

		So(len(config.Hosts), ShouldEqual, 0)
		So(config.Defaults.Port, ShouldEqual, uint(0))
		So(config.Defaults.HostName, ShouldEqual, "")
		So(config.Defaults.User, ShouldEqual, "")
	})
}

func dummyConfig() *Config {
	config := New()
	config.Hosts["toto"] = Host{
		HostName: "1.2.3.4",
	}
	config.Hosts["titi"] = Host{
		HostName:     "tata",
		Port:         23,
		User:         "moul",
		ProxyCommand: "nc -v 4242",
	}
	config.Hosts["tonton"] = Host{
		ResolveNameservers: []string{"a.com", "1.2.3.4"},
	}
	config.Hosts["toutou"] = Host{
		ResolveCommand: "dig -t %h",
	}
	config.Hosts["tutu"] = Host{
		Gateways: []string{"titi", "direct", "1.2.3.4"},
		Inherits: []string{"toto", "tutu", "*.ddd"},
	}
	config.Hosts["empty"] = Host{}
	config.Hosts["tata"] = Host{
		Inherits: []string{"tutu", "titi", "toto", "tutu"},
	}
	config.Hosts["*.ddd"] = Host{
		HostName:               "1.3.5.7",
		PasswordAuthentication: "yes",
	}
	config.Defaults = Host{
		Port: 22,
		User: "root",
	}
	config.applyMissingNames()
	return config
}

func TestConfig(t *testing.T) {
	Convey("Testing dummyConfig", t, func() {
		config := dummyConfig()

		So(len(config.Hosts), ShouldEqual, 8)

		So(config.Hosts["toto"].HostName, ShouldEqual, "1.2.3.4")
		So(config.Hosts["toto"].Port, ShouldEqual, 0)
		So(config.Hosts["toto"].name, ShouldEqual, "toto")
		So(config.Hosts["toto"].isDefault, ShouldEqual, false)

		So(config.Hosts["titi"].HostName, ShouldEqual, "tata")
		So(config.Hosts["titi"].User, ShouldEqual, "moul")
		So(config.Hosts["titi"].ProxyCommand, ShouldEqual, "nc -v 4242")
		So(config.Hosts["titi"].Port, ShouldEqual, 23)
		So(config.Hosts["titi"].isDefault, ShouldEqual, false)

		So(config.Hosts["tonton"].isDefault, ShouldEqual, false)
		So(config.Hosts["tonton"].Port, ShouldEqual, uint(0))
		So(config.Hosts["tonton"].ResolveNameservers, ShouldResemble, []string{"a.com", "1.2.3.4"})

		So(config.Hosts["toutou"].isDefault, ShouldEqual, false)
		So(config.Hosts["toutou"].Port, ShouldEqual, uint(0))
		So(config.Hosts["toutou"].ResolveCommand, ShouldEqual, "dig -t %h")

		So(config.Hosts["tutu"].isDefault, ShouldEqual, false)
		So(config.Hosts["tutu"].Port, ShouldEqual, uint(0))
		So(config.Hosts["tutu"].Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})

		So(config.Hosts["*.ddd"].isDefault, ShouldEqual, false)
		So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")

		So(config.Hosts["empty"].isDefault, ShouldEqual, false)
		So(config.Hosts["empty"].Port, ShouldEqual, uint(0))

		So(config.Defaults.User, ShouldEqual, "root")
		So(config.Defaults.Port, ShouldEqual, uint(22))
		So(config.Defaults.isDefault, ShouldEqual, true)
	})
}

func TestConfig_LoadConfig(t *testing.T) {
	Convey("Testing Config.LoadConfig", t, func() {

		config := New()
		err := config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)
		So(len(config.Hosts), ShouldEqual, 4)
		So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
		So(config.Hosts["aaa"].Port, ShouldEqual, uint(0))
		So(config.Hosts["aaa"].User, ShouldEqual, "")
		So(config.Hosts["bbb"].HostName, ShouldEqual, "")
		So(config.Hosts["bbb"].Port, ShouldEqual, uint(21))
		So(config.Hosts["bbb"].User, ShouldEqual, "")
		So(config.Hosts["ccc"].HostName, ShouldEqual, "5.6.7.8")
		So(config.Hosts["ccc"].Port, ShouldEqual, uint(24))
		So(config.Hosts["ccc"].User, ShouldEqual, "toor")
		So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")
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
      "HostName": "1.3.5.7",
      "PasswordAuthentication": "yes"
    },
    "empty": {},
    "tata": {
      "Inherits": [
        "tutu",
        "titi",
        "toto",
        "tutu"
      ]
    },
    "titi": {
      "HostName": "tata",
      "Port": 23,
      "User": "moul",
      "ProxyCommand": "nc -v 4242"
    },
    "tonton": {
      "ResolveNameservers": [
        "a.com",
        "1.2.3.4"
      ]
    },
    "toto": {
      "HostName": "1.2.3.4"
    },
    "toutou": {
      "ResolveCommand": "dig -t %h"
    },
    "tutu": {
      "Inherits": [
        "toto",
        "tutu",
        "*.ddd"
      ],
      "Gateways": [
        "titi",
        "direct",
        "1.2.3.4"
      ]
    }
  },
  "defaults": {
    "Port": 22,
    "User": "root"
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
			host, err = config.getHostByName("titi", false, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")

			host, err = config.getHostByName("titi", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")

			host, err = config.getHostByName("dontexists", false, true)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("dontexists", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists")

			host, err = config.getHostByName("regex.ddd", false, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")

			host, err = config.getHostByName("regex.ddd", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
		})

		Convey("With gateway", func() {
			host, err = config.getHostByName("titi/gateway", false, true)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("titi/gateway", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByName("dontexists/gateway", false, true)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("dontexists/gateway", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByName("regex.ddd/gateway", false, true)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("regex.ddd/gateway", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd/gateway")
			So(host.HostName, ShouldNotEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})
	})
}

func TestConfig_GetGatewaySafe(t *testing.T) {
	Convey("Testing Config.GetGatewaySafe", t, func() {

		config := dummyConfig()
		var host *Host

		Convey("Without gateway", func() {
			host = config.GetGatewaySafe("titi")
			So(host.Name(), ShouldEqual, "titi")

			host = config.GetGatewaySafe("dontexists")
			So(host.Name(), ShouldEqual, "dontexists")

			host = config.GetGatewaySafe("regex.ddd")
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
		})

		Convey("With gateway", func() {
			host = config.GetGatewaySafe("titi/gateway")
			So(host.Name(), ShouldEqual, "titi/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetGatewaySafe("dontexists/gateway")
			So(host.Name(), ShouldEqual, "dontexists/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetGatewaySafe("regex.ddd/gateway")
			So(host.Name(), ShouldEqual, "regex.ddd/gateway")
			So(host.HostName, ShouldNotEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})
	})
}

func TestConfig_LoadFiles(t *testing.T) {
	Convey("Testing Config.LoadFiles", t, func() {
		config := New()
		file, err := ioutil.TempFile(os.TempDir(), "assh-tests")
		So(err, ShouldBeNil)
		defer os.Remove(file.Name())
		file.Write([]byte(configExample))

		Convey("Loading a simple file", func() {
			err = config.LoadFiles(file.Name())

			So(err, ShouldBeNil)
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 1)
			So(len(config.Hosts), ShouldEqual, 4)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, uint(0))
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "")
			So(config.Hosts["bbb"].Port, ShouldEqual, uint(21))
			So(config.Hosts["bbb"].User, ShouldEqual, "")
			So(config.Hosts["ccc"].HostName, ShouldEqual, "5.6.7.8")
			So(config.Hosts["ccc"].Port, ShouldEqual, uint(24))
			So(config.Hosts["ccc"].User, ShouldEqual, "toor")
			So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")
			So(config.Hosts["*.ddd"].Port, ShouldEqual, uint(0))
			So(config.Hosts["*.ddd"].User, ShouldEqual, "")
			So(config.Defaults.Port, ShouldEqual, uint(22))
			So(config.Defaults.User, ShouldEqual, "root")
		})
		Convey("Loading the same file again", func() {
			config.LoadFiles(file.Name())
			err = config.LoadFiles(file.Name())

			So(err, ShouldBeNil)
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 1)
			So(len(config.Hosts), ShouldEqual, 4)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, uint(0))
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "")
			So(config.Hosts["bbb"].Port, ShouldEqual, uint(21))
			So(config.Hosts["bbb"].User, ShouldEqual, "")
			So(config.Hosts["ccc"].HostName, ShouldEqual, "5.6.7.8")
			So(config.Hosts["ccc"].Port, ShouldEqual, uint(24))
			So(config.Hosts["ccc"].User, ShouldEqual, "toor")
			So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")
			So(config.Hosts["*.ddd"].Port, ShouldEqual, uint(0))
			So(config.Hosts["*.ddd"].User, ShouldEqual, "")
			So(config.Defaults.Port, ShouldEqual, uint(22))
			So(config.Defaults.User, ShouldEqual, "root")
		})
	})
	// FIXME: test globbing
}

func TestConfig_getHostByPath(t *testing.T) {
	Convey("Testing Config.getHostByPath", t, func() {

		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.getHostByPath("titi", false, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("titi", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("dontexists", false, true)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByPath("dontexists", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("regex.ddd", false, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("regex.ddd", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host, err = config.getHostByPath("titi/gateway", false, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("titi/gateway", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("dontexists/gateway", false, true)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByPath("dontexists/gateway", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("regex.ddd/gateway", false, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("regex.ddd/gateway", true, true)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
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
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.GetHost("dontexists")
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.GetHost("regex.ddd")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host, err = config.GetHost("titi/gateway")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.GetHost("dontexists/gateway")
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			// FIXME: check if this is a normal behavior
			host, err = config.GetHost("regex.ddd/gateway")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})

		Convey("Inheritance", FailureContinues, func() {
			host, err = config.GetHost("tata")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"tata": true,
				"tutu": true,
				"titi": true,
				"toto": true,
			})
			So(host.ProxyCommand, ShouldEqual, "nc -v 4242")
			So(host.User, ShouldEqual, "moul")
			So(host.Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})
			So(host.PasswordAuthentication, ShouldEqual, "yes")

			host, err = config.GetHost("tutu")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"tutu":  true,
				"toto":  true,
				"*.ddd": true,
			})
			So(host.User, ShouldEqual, "root")
			So(host.Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})
			So(host.PasswordAuthentication, ShouldEqual, "yes")
		})
	})
}

func TestConfig_GetHostSafe(t *testing.T) {
	Convey("Testing Config.GetHostSafe", t, func() {

		config := dummyConfig()
		var host *Host

		Convey("Without gateway", func() {
			host = config.GetHostSafe("titi")
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetHostSafe("dontexists")
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetHostSafe("regex.ddd")
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host = config.GetHostSafe("titi/gateway")
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host = config.GetHostSafe("dontexists/gateway")
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 1)

			host = config.GetHostSafe("regex.ddd/gateway")
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})
	})
}

func TestConfig_WriteSshConfig(t *testing.T) {
	Convey("Testing Config.WriteSshConfig", t, func() {
		config := dummyConfig()

		var buffer bytes.Buffer

		err := config.WriteSshConfigTo(&buffer)
		So(err, ShouldBeNil)

		expected := `# ssh config generated by advanced-ssh-config

# host-based configuration
Host *.ddd
  HostName 1.3.5.7
  PasswordAuthentication yes

Host empty

Host tata
  HostName 1.2.3.4
  PasswordAuthentication yes
  Port 22
  User moul
#  ProxyCommand nc -v 4242

Host titi
  HostName tata
  Port 23
  User moul
#  ProxyCommand nc -v 4242

Host tonton

Host toto
  HostName 1.2.3.4

Host toutou

Host tutu
  HostName 1.2.3.4
  PasswordAuthentication yes
  Port 22

# global configuration
Host *
  Port 22
  User root
  ProxyCommand assh proxy --port=%p %h
`
		So(buffer.String(), ShouldEqual, expected)
	})
}
