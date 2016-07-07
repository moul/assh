package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHostsListToList(t *testing.T) {
	Convey("Testing HostsList.ToList()", t, func() {
		m := HostsMap{
			"aaa": &Host{name: "aaa"},
			"bbb": &Host{name: "bbb"},
			"ccc": &Host{name: "ccc"},
		}

		list := m.ToList()
		So(len(list), ShouldEqual, 3)
	})
}

func TestHostsListSortedList(t *testing.T) {
	Convey("Testing HostsList.SortedList()", t, func() {
		m := HostsMap{
			"ccc": &Host{name: "ccc"},
			"ddd": &Host{name: "ddd"},
			"aaa": &Host{name: "aaa"},
			"bbb": &Host{name: "bbb"},
		}

		sorted := m.SortedList()

		So(sorted[0].name, ShouldEqual, "aaa")
		So(sorted[1].name, ShouldEqual, "bbb")
		So(sorted[2].name, ShouldEqual, "ccc")
		So(sorted[3].name, ShouldEqual, "ddd")
	})
}
