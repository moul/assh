package commands

import (
	"strings"
	"testing"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/vendor/github.com/smartystreets/goconvey/convey"
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
defaults:
  Port: 22
  User: root
`

func TestComputeHost(t *testing.T) {
	Convey("Testing computeHost()", t, func() {
		config := config.New()

		err := config.LoadConfig(strings.NewReader(configExample))
		host, err := computeHost("aaa", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "1.2.3.4")
		So(host.Port, ShouldEqual, 22)

		err = config.LoadConfig(strings.NewReader(configExample))
		host, err = computeHost("aaa", 42, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "1.2.3.4")
		So(host.Port, ShouldEqual, 42)

		err = config.LoadConfig(strings.NewReader(configExample))
		host, err = computeHost("eee", 0, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "eee")
		So(host.Port, ShouldEqual, 22)

		err = config.LoadConfig(strings.NewReader(configExample))
		host, err = computeHost("eee", 42, config)
		So(err, ShouldBeNil)
		So(host.HostName, ShouldEqual, "eee")
		So(host.Port, ShouldEqual, 42)
	})
}

func TestCommandApplyHost(t *testing.T) {
	Convey("Testing commandApplyHost()", t, func() {
		host := config.NewHost("abc")
		host.HostName = "1.2.3.4"
		host.Port = 42

		var input, output, expected string

		input = "ls -la"
		output = commandApplyHost(input, host)
		expected = "ls -la"
		So(output, ShouldEqual, expected)

		input = "nc %h %p"
		output = commandApplyHost(input, host)
		expected = "nc 1.2.3.4 42"
		So(output, ShouldEqual, expected)

		input = "ssh %name"
		output = commandApplyHost(input, host)
		expected = "ssh abc"
		So(output, ShouldEqual, expected)

		input = "echo %h %p %name %h %p %name"
		output = commandApplyHost(input, host)
		expected = "echo 1.2.3.4 42 abc 1.2.3.4 42 abc"
		So(output, ShouldEqual, expected)
	})
}
