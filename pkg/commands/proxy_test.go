package commands

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"moul.io/assh/v2/pkg/config"
)

const configExample string = `
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
  eee:
    ResolveCommand: /bin/sh -c "echo 42.42.42.42"
defaults:
  Port: 22
  User: root
`

func TestComputeHost(t *testing.T) {
	Convey("Testing computeHost()", t, func() {
		config := config.New()

		err := config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)
		host, err := computeHost("aaa", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "1.2.3.4")
		So(host.Port, ShouldEqual, "22")

		err = config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)
		host, err = computeHost("aaa", 42, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "1.2.3.4")
		So(host.Port, ShouldEqual, "42")

		err = config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)
		host, err = computeHost("eee", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "eee")
		So(host.Port, ShouldEqual, "22")

		err = config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)
		host, err = computeHost("eee", 42, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "eee")
		So(host.Port, ShouldEqual, "42")
	})
}

func Test_runProxy(t *testing.T) {
	Convey("Testing proxyCommand()", t, func() {
		// FIXME: test stdout
		config := config.New()
		err := config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)
		host, err := computeHost("aaa", 0, config)
		So(err, ShouldBeNil)

		err = runProxy(host, "echo test from proxyCommand", false)
		So(err, ShouldBeNil)

		err = runProxy(host, "/bin/sh -c 'echo test from proxyCommand'", false)
		So(err, ShouldBeNil)

		err = runProxy(host, "/bin/sh -c 'exit 1'", false)
		So(err, ShouldNotBeNil)

		err = runProxy(host, "blah", true)
		So(err, ShouldResemble, fmt.Errorf("dry-run: Execute [blah]"))
	})
}

func Test_hostPrepare(t *testing.T) {
	Convey("Testing hostPrepare()", t, func() {
		config := config.New()
		err := config.LoadConfig(strings.NewReader(configExample))
		So(err, ShouldBeNil)

		host, err := computeHost("aaa", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "1.2.3.4")
		So(hostPrepare(host, ""), ShouldBeNil)
		So(host.HostName, ShouldEqual, "1.2.3.4")

		host, err = computeHost("bbb", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "bbb")
		So(hostPrepare(host, ""), ShouldBeNil)
		So(host.HostName, ShouldEqual, "bbb")

		host, err = computeHost("eee", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "eee")
		So(hostPrepare(host, ""), ShouldBeNil)
		So(host.HostName, ShouldEqual, "42.42.42.42")
	})
}
