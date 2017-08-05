package logging

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestLoggingWithNilReferenceProducesTraditionalBehavior(t *testing.T) {
	out := new(bytes.Buffer)
	log.SetOutput(out)
	log.SetFlags(0)
	defer log.SetOutput(os.Stdout)
	defer log.SetFlags(log.LstdFlags)

	thing := new(ThingUnderTest)
	thing.Action()

	assertions.New(t).So(out.String(), should.Equal, "Hello, World!\n")
}

func TestLoggingWithLoggerCapturesOutput(t *testing.T) {
	out := new(bytes.Buffer)
	log.SetOutput(out)
	log.SetFlags(0)
	defer log.SetOutput(os.Stdout)
	defer log.SetFlags(log.LstdFlags)

	thing := new(ThingUnderTest)
	thing.log = Capture()
	thing.Action()

	assertions.New(t).So(thing.log.Log.String(), should.Equal, "Hello, World!\n")
	assertions.New(t).So(out.Len(), should.Equal, 0)
}

func TestLogCallsAreCounted(t *testing.T) {
	thing := new(ThingUnderTest)
	thing.log = Capture()
	for x := 0; x < 10; x++ {
		thing.Action()
	}
	assertions.New(t).So(thing.log.Calls, should.Equal, 10)
}

func TestLoggingWithDiscard(t *testing.T) {
	out := new(bytes.Buffer)
	log.SetOutput(out)
	log.SetFlags(0)
	defer log.SetOutput(os.Stdout)
	defer log.SetFlags(log.LstdFlags)

	thing := new(ThingUnderTest)
	thing.log = Discard()
	thing.Action()

	assertions.New(t).So(thing.log.Log.Len(), should.Equal, 0)
	assertions.New(t).So(out.Len(), should.Equal, 0)
}

/////////////////////////////////////////////////////////////////////////////////

type ThingUnderTest struct {
	log *Logger
}

func (this *ThingUnderTest) Action() {
	this.log.Printf("Hello, World!")
}
