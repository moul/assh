package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptionString(t *testing.T) {
	Convey("Testing Option.String()", t, func() {
		option := Option{
			Name:  "name",
			Value: "value",
		}
		So(option.String(), ShouldEqual, "name=value")
	})
}

func TestOptionsListToStringList(t *testing.T) {
	Convey("Testing OptionsList.ToStringList()", t, func() {
		ol := OptionsList{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
			{Name: "name3", Value: "value3"},
		}
		So(ol.ToStringList(), ShouldResemble, []string{"name1=value1", "name2=value2", "name3=value3"})
	})
}

func TestOptionsListRemove(t *testing.T) {
	Convey("Testing OptionsList.Remove()", t, func() {
		ol := OptionsList{
			{Name: "name1", Value: "value1"},
			{Name: "name2", Value: "value2"},
			{Name: "name3", Value: "value3"},
		}
		So(ol.ToStringList(), ShouldResemble, []string{"name1=value1", "name2=value2", "name3=value3"})

		ol.Remove("name4")
		So(ol.ToStringList(), ShouldResemble, []string{"name1=value1", "name2=value2", "name3=value3"})

		ol.Remove("name2")
		So(ol.ToStringList(), ShouldResemble, []string{"name1=value1", "name3=value3"})

		ol.Remove("name2")
		So(ol.ToStringList(), ShouldResemble, []string{"name1=value1", "name3=value3"})

		ol.Remove("name3")
		So(ol.ToStringList(), ShouldResemble, []string{"name1=value1"})

		ol.Remove("name1")
		So(ol.ToStringList(), ShouldResemble, []string{})

		ol.Remove("name1")
		So(ol.ToStringList(), ShouldResemble, []string{})
	})
}
