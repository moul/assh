# logging
--
    import "github.com/smartystreets/logging"

package logging implements a Logger that, when nil, forwards to the
corresponding functions in the standard log package. When not nil, it captures
log calls in a buffer for later inspection. This can be useful when needing to
inspect or squelch log output from test code. The main advantage to this
approach is that it is not necessary to provide a non-nil instance in
'contructor' functions or wireup for production code. It is also still trivial
to set a non-nil reference in test code.

## Usage

#### type Logger

```go
type Logger struct {
	Log *bytes.Buffer
	*log.Logger
}
```

Logger is meant be included as a pointer field on a struct. Leaving the instance
as a nil reference will cause any calls on the *Logger to forward to the
corresponding functions from the standard log package. This is meant to be the
behavior in production. In testing, set the field to a non-nil instance of a
*Logger to record log statements for later inspection.

#### func  Capture

```go
func Capture() *Logger
```
Capture creates a new *Logger instance with an internal buffer. The prefix and
flags default to the values of log.Prefix() and log.Flags(), respectively. This
function is meant to be called from test code. See the godoc for the Logger
struct for details.

#### func (*Logger) Fatal

```go
func (this *Logger) Fatal(v ...interface{})
```
Fatal -> log.Fatal

#### func (*Logger) Fatalf

```go
func (this *Logger) Fatalf(format string, v ...interface{})
```
Fatalf -> log.Fatalf

#### func (*Logger) Fatalln

```go
func (this *Logger) Fatalln(v ...interface{})
```
Fatalln -> log.Fatalln

#### func (*Logger) Flags

```go
func (this *Logger) Flags() int
```
Flags -> log.Flags

#### func (*Logger) Output

```go
func (this *Logger) Output(calldepth int, s string) error
```
Output -> log.Output

#### func (*Logger) Panic

```go
func (this *Logger) Panic(v ...interface{})
```
Panic -> log.Panic

#### func (*Logger) Panicf

```go
func (this *Logger) Panicf(format string, v ...interface{})
```
Panicf -> log.Panicf

#### func (*Logger) Panicln

```go
func (this *Logger) Panicln(v ...interface{})
```
Panicln -> log.Panicln

#### func (*Logger) Prefix

```go
func (this *Logger) Prefix() string
```
Prefix -> log.Prefix

#### func (*Logger) Print

```go
func (this *Logger) Print(v ...interface{})
```
Print -> log.Print

#### func (*Logger) Printf

```go
func (this *Logger) Printf(format string, v ...interface{})
```
Printf -> log.Printf

#### func (*Logger) Println

```go
func (this *Logger) Println(v ...interface{})
```
Println -> log.Println

#### func (*Logger) SetFlags

```go
func (this *Logger) SetFlags(flag int)
```
SetFlags -> log.SetFlags

#### func (*Logger) SetOutput

```go
func (this *Logger) SetOutput(w io.Writer)
```
SetOutput -> log.SetOutput

#### func (*Logger) SetPrefix

```go
func (this *Logger) SetPrefix(prefix string)
```
SetPrefix -> log.SetPrefix
