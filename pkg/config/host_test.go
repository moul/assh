package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHost_ApplyDefaults(t *testing.T) {
	Convey("Testing Host.ApplyDefaults", t, func() {
		Convey("Standard configuration", func() {
			host := &Host{
				Name: "example",
				Host: "example.com",
				User: "root",
			}
			defaults := &Host{
				User: "bobby",
				Port: 42,
			}
			host.ApplyDefaults(defaults)
			So(host.Port, ShouldEqual, uint(42))
			So(host.Name, ShouldEqual, "example")
			So(host.Host, ShouldEqual, "example.com")
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
			So(host.Port, ShouldEqual, uint(22))
			So(host.Name, ShouldEqual, "")
			So(host.Host, ShouldEqual, "")
			So(host.User, ShouldEqual, "")
			So(len(host.Gateways), ShouldEqual, 0)
			So(host.ProxyCommand, ShouldEqual, "")
			So(len(host.ResolveNameservers), ShouldEqual, 0)
			So(host.ResolveCommand, ShouldEqual, "")
			So(host.ControlPath, ShouldEqual, "")
		})
	})
}
