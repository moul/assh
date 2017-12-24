package flexyaml

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMakeFlexible(t *testing.T) {
	Convey("Testing MakeFlexible()", t, FailureContinues, func() {
		var input, output []byte
		var err error

		input = []byte("somekey: ok")
		output, err = MakeFlexible(input)
		So(err, ShouldBeNil)
		So(output, ShouldResemble, []byte("somekey: ok"))

		input = []byte("SomeKey: ok")
		output, err = MakeFlexible(input)
		So(err, ShouldBeNil)
		So(output, ShouldResemble, []byte("somekey: ok"))

		input = []byte("SOMEKEY: ok")
		output, err = MakeFlexible(input)
		So(err, ShouldBeNil)
		So(output, ShouldResemble, []byte("somekey: ok"))

		input = []byte("sOmEkEy: ok")
		output, err = MakeFlexible(input)
		So(err, ShouldBeNil)
		So(output, ShouldResemble, []byte("somekey: ok"))

		expected := []byte(`hosts:

  aaa:
    hostname: 1.2.3.4

  bbb:
    port: 21

  ccc:
    hostname: 5.6.7.8
    port: 24
    user: Toor
  "*.ddd":
    hostname: 1.3.5.7

  eee:
    inherits:
    - aaa
    - bBb
    - aaa

  fff:
    inherits:
    - bbb
    - eee
    - "*.ddd"

  ggg:
    gateways:
    - direct
    - fff

  hhh:
    gateways:
    - ggg
    - direct

  iii:
    gateways:
    - test.ddd

  jjj:
    hostname: "%h.jjjjj"

  "*.kkk":
    hostname: "%h.kkkkk"

defaults:
  port: 22
  user: root

includes:
  - /path/to/dir/*.yml
  - /path/to/file.yml
`)
		input = []byte(`hosts:

  aaa:
    HostName: 1.2.3.4

  bBb:
    Port: 21

  ccc:
    HostName: 5.6.7.8
    Port: 24
    USER: Toor
  "*.ddd":
    hostName: 1.3.5.7

  eee:
    Inherits:
    - aaa
    - bBb
    - aaa

  fff:
    Inherits:
    - bbb
    - eee
    - "*.ddd"

  ggg:
    Gateways:
    - direct
    - fff

  hhh:
    Gateways:
    - ggg
    - direct

  iii:
    Gateways:
    - test.ddd

  jjj:
    HostName: "%h.jjjjj"

  "*.kkk":
    HostName: "%h.kkkkk"

defaults:
  Port: 22
  User: root

includes:
  - /path/to/dir/*.yml
  - /path/to/file.yml
`)
		output, err = MakeFlexible(input)
		So(err, ShouldBeNil)
		So(string(output), ShouldEqual, string(expected))
	})
}

func TestUnmarshal(t *testing.T) {
	Convey("Testing Unmarshal()", t, FailureContinues, func() {
		type C struct {
			SomeKey string
		}
		var out C

		err := Unmarshal([]byte("somekey: ok"), &out)
		So(err, ShouldBeNil)
		So(out.SomeKey, ShouldEqual, "ok")
		out = C{}

		err = Unmarshal([]byte("SomeKey: ok"), &out)
		So(err, ShouldBeNil)
		So(out.SomeKey, ShouldEqual, "ok")
		out = C{}

		err = Unmarshal([]byte("SOMEKEY: ok"), &out)
		So(err, ShouldBeNil)
		So(out.SomeKey, ShouldEqual, "ok")
		out = C{}

		err = Unmarshal([]byte("sOmEkEy: ok"), &out)
		So(err, ShouldBeNil)
		So(out.SomeKey, ShouldEqual, "ok")
	})
}
