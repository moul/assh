package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHost_ApplyDefaults(t *testing.T) {
	Convey("Testing Host.ApplyDefaults", t, func() {
		Convey("Standard configuration", func() {
			host := &Host{
				name:     "example",
				HostName: "example.com",
				User:     "root",
			}
			defaults := &Host{
				User: "bobby",
				Port: "42",
			}
			host.ApplyDefaults(defaults)
			So(host.Port, ShouldEqual, "42")
			So(host.Name(), ShouldEqual, "example")
			So(host.HostName, ShouldEqual, "example.com")
			So(host.User, ShouldEqual, "root")
			So(len(host.Gateways), ShouldEqual, 0)
			So(host.ProxyCommand, ShouldEqual, "")
			So(len(host.ResolveNameservers), ShouldEqual, 0)
			So(host.ResolveCommand, ShouldEqual, "")
			So(host.ControlPath, ShouldEqual, "")
		})
		Convey("Empty configuration", func() {
			host := &Host{}
			defaults := &Host{}
			host.ApplyDefaults(defaults)
			So(host.Port, ShouldEqual, "22")
			So(host.Name(), ShouldEqual, "")
			So(host.HostName, ShouldEqual, "")
			So(host.User, ShouldEqual, "")
			So(len(host.Gateways), ShouldEqual, 0)
			So(host.ProxyCommand, ShouldEqual, "")
			So(len(host.ResolveNameservers), ShouldEqual, 0)
			So(host.ResolveCommand, ShouldEqual, "")
			So(host.ControlPath, ShouldEqual, "")
		})
	})
}

func TestHost_ExpandString(t *testing.T) {
	Convey("Testing Host.ExpandString()", t, func() {
		host := NewHost("abc")
		host.HostName = "1.2.3.4"
		host.Port = "42"

		var input, output, expected string

		input = "ls -la"
		output = host.ExpandString(input)
		expected = "ls -la"
		So(output, ShouldEqual, expected)

		input = "nc %h %p"
		output = host.ExpandString(input)
		expected = "nc 1.2.3.4 42"
		So(output, ShouldEqual, expected)

		input = "ssh %name"
		output = host.ExpandString(input)
		expected = "ssh abc"
		So(output, ShouldEqual, expected)

		input = "echo %h %p %name %h %p %name"
		output = host.ExpandString(input)
		expected = "echo 1.2.3.4 42 abc 1.2.3.4 42 abc"
		So(output, ShouldEqual, expected)
	})
}
