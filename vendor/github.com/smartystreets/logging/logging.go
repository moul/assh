// package logging implements a Logger that, when nil, forwards to the
// corresponding functions in the standard log package. When not nil,
// it captures log calls in a buffer for later inspection. This can be
// useful when needing to inspect or squelch log output from test code.
// The main advantage to this approach is that it is not necessary to
// provide a non-nil instance in 'constructor' functions or wireup for
// production code. It is also still trivial to set a non-nil reference
// in test code.
// This library requires go 1.5+
package logging

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Logger is meant be included as a pointer field on a struct. Leaving the
// instance as a nil reference will cause any calls on the *Logger to forward
// to the corresponding functions from the standard log package. This is meant
// to be the behavior in production. In testing, set the field to a non-nil
// instance of a *Logger to record log statements for later inspection.
type Logger struct {
	*log.Logger

	Log   *bytes.Buffer
	Calls int
}

// Capture creates a new *Logger instance with an internal buffer. The prefix
// and flags default to the values of log.Prefix() and log.Flags(), respectively.
// This function is meant to be called from test code. See the godoc for the
// Logger struct for details.
func Capture() *Logger {
	out := new(bytes.Buffer)
	return &Logger{
		Log:    out,
		Logger: log.New(out, log.Prefix(), log.Flags()),
	}
}

// Discard creates a new *Logger instance with its internal buffer set to
// ioutil.Discard. This is useful if you want your production code to be
// quiet but your test code to be verbose. In that case, use Discard()
// in production code and Capture() in test code.
func Discard() *Logger {
	return &Logger{
		Log:    new(bytes.Buffer),
		Logger: log.New(ioutil.Discard, "", 0),
	}
}

// SetOutput -> log.SetOutput
func (this *Logger) SetOutput(w io.Writer) {
	if this == nil {
		log.SetOutput(w)
	} else {
		this.Logger.SetOutput(w)
	}
}

// Output -> log.Output
func (this *Logger) Output(calldepth int, s string) error {
	if this == nil {
		return log.Output(calldepth, s)
	}
	this.Calls++
	return this.Logger.Output(calldepth, s)
}

// Fatal -> log.Fatal
func (this *Logger) Fatal(v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprint(v...))
		os.Exit(1)
	} else {
		this.Calls++
		this.Logger.Fatal(v...)
	}
}

// Fatalf -> log.Fatalf
func (this *Logger) Fatalf(format string, v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprintf(format, v...))
		os.Exit(1)
	} else {
		this.Calls++
		this.Logger.Fatalf(format, v...)
	}
}

// Fatalln -> log.Fatalln
func (this *Logger) Fatalln(v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprintln(v...))
		os.Exit(1)
	} else {
		this.Calls++
		this.Logger.Fatalln(v...)
	}
}

// Flags -> log.Flags
func (this *Logger) Flags() int {
	if this == nil {
		return log.Flags()
	}
	return this.Logger.Flags()
}

// Panic -> log.Panic
func (this *Logger) Panic(v ...interface{}) {
	if this == nil {
		s := fmt.Sprint(v...)
		this.Output(3, s)
		panic(s)
	} else {
		this.Calls++
		this.Logger.Panic(v...)
	}
}

// Panicf -> log.Panicf
func (this *Logger) Panicf(format string, v ...interface{}) {
	if this == nil {
		s := fmt.Sprintf(format, v...)
		this.Output(3, s)
		panic(s)
	} else {
		this.Calls++
		this.Logger.Panicf(format, v...)
	}
}

// Panicln -> log.Panicln
func (this *Logger) Panicln(v ...interface{}) {
	if this == nil {
		s := fmt.Sprintln(v...)
		this.Output(3, s)
		panic(s)
	} else {
		this.Calls++
		this.Logger.Panicln(v...)
	}
}

// Prefix -> log.Prefix
func (this *Logger) Prefix() string {
	if this == nil {
		return log.Prefix()
	}
	this.Calls++
	return this.Logger.Prefix()
}

// Print -> log.Print
func (this *Logger) Print(v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprint(v...))
	} else {
		this.Calls++
		this.Logger.Print(v...)
	}
}

// Printf -> log.Printf
func (this *Logger) Printf(format string, v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprintf(format, v...))
	} else {
		this.Calls++
		this.Logger.Printf(format, v...)
	}
}

// Println -> log.Println
func (this *Logger) Println(v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprintln(v...))
	} else {
		this.Calls++
		this.Logger.Println(v...)
	}
}

// SetFlags -> log.SetFlags
func (this *Logger) SetFlags(flag int) {
	if this == nil {
		log.SetFlags(flag)
	} else {
		this.Logger.SetFlags(flag)
	}
}

// SetPrefix -> log.SetPrefix
func (this *Logger) SetPrefix(prefix string) {
	if this == nil {
		log.SetPrefix(prefix)
	} else {
		this.Logger.SetPrefix(prefix)
	}
}
