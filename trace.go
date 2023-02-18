package oops

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/calebcase/oops/lines"
)

// TraceSkipInternal is the number of frames created by internal oops calls.
// This is the number of frames to skip if you would like to only include
// frames up til the site where TraceN is called.
const TraceSkipInternal = 7

// Capturer provides a method for capturing trace data.
type Capturer interface {
	Capture(err error, skip int) (data any)
}

// CaptureFunc defines a Capturer that calls the given function to generate the
// capture.
type CaptureFunc[T any] func(error, int) T

// Capture implements Capturer.
func (cf CaptureFunc[T]) Capture(err error, skip int) any {
	return cf(err, skip)
}

// Frames wraps []runtime.Frame to provide a custom stringer.
type Frames []runtime.Frame

// String returns the frames formatted as an numbered intended array of frames.
func (fs Frames) String() string {
	ls := []string{}

	for i, f := range fs {
		prefix := fmt.Sprintf("[%d] ", i)
		ls = append(ls, prefix+f.Function)
		ls = append(ls, strings.Repeat(" ", len(prefix))+fmt.Sprintf("%s:%d", f.File, f.Line))
	}

	return strings.Join(ls, "\n")
}

// CaptureRuntimeFrames returns the captured stack as []runtime.Frame.
func CaptureRuntimeFrames(_ error, skip int) []runtime.Frame {
	callers := make([]uintptr, 10)
	n := runtime.Callers(skip, callers)
	callers = callers[:n]

	cfs := runtime.CallersFrames(callers)

	fs := make([]runtime.Frame, 0, n)
	for {
		f, more := cfs.Next()
		if !more {
			break
		}

		fs = append(fs, f)
	}

	return fs
}

// CaptureFrames returns the captured stack as Frames.
func CaptureFrames(err error, skip int) (data Frames) {
	return Frames(CaptureRuntimeFrames(err, skip))
}

// defaultCapturer is the package level setting for the capture function. This
// is what will be used if no trace options are provided.
var defaultCapturer Capturer = CaptureFunc[Frames](CaptureFrames)

// SetDefaultCapturer changes the default trace capturer used by calls to
// Trace, TraceN, and TraceWithOptions. The default capturer creates a capture
// using CaptureFrames.
func SetDefaultCapturer(c Capturer) {
	defaultCapturer = c
}

// TraceError is an error with trace data.
type TraceError struct {
	Data any
	Err  error
}

// Error implements error.
func (te *TraceError) Error() string {
	return fmt.Sprintf("%v", te)
}

// Unwrap implements the implied interface for errors.Unwrap.
func (te *TraceError) Unwrap() error {
	if te == nil || te.Err == nil {
		return nil
	}

	return te.Err
}

// Format implements fmt.Format.
func (te *TraceError) Format(f fmt.State, verb rune) {
	if te == nil || te.Err == nil {
		fmt.Fprintf(f, "<nil>")

		return
	}

	flag := ""
	if f.Flag(int('+')) {
		flag = "+"
	}

	if flag == "" {
		fmt.Fprintf(f, "%"+string(verb), te.Err)

		return
	}

	output := []string{}
	output = append(output, lines.Indent(lines.Sprintf("%"+flag+string(verb), te.Err), "··", 1)...)
	output = append(output, lines.Indent(lines.Sprintf("%"+flag+string(verb), te.Data), "··", 0)...)

	f.Write([]byte(strings.Join(output, "\n")))
}

// MarshalJSON implements json.Marshaler.
func (te *TraceError) MarshalJSON() (bs []byte, err error) {
	if te == nil || te.Err == nil {
		return []byte("null"), nil
	}

	ebs, err := ErrorMarshalJSON(te.Err)
	if err != nil {
		return nil, err
	}

	output := struct {
		Type string          `json:"type"`
		Err  json.RawMessage `json:"err"`
		Data any             `json:"data"`
	}{
		Type: fmt.Sprintf("%T", te.Err),
		Err:  json.RawMessage(ebs),
		Data: te.Data,
	}

	return json.Marshal(output)
}

// Trace captures a trace and combines it with err. Tracer is use to capture
// the trace data and internal stacks are skipped.
func Trace(err error) error {
	return TraceN(err, TraceSkipInternal)
}

// TraceN captures a trace and combines it with err. Capturer is use to capture
// the trace data with skipped levels.
func TraceN(err error, skip int) error {
	return traceWithOptions(err, TraceOptions{
		Skip: skip,
	})
}

// TraceOptions allows setting specific options for the trace.
type TraceOptions struct {
	// Skip the given number of levels of tracing. Default is 0.
	Skip int

	// Capturer controls the specific capturer implementation to use. If
	// not set, then the package level Capture will be used.
	Capturer Capturer
}

// TraceWithOptions captures a trace using the given options.
func TraceWithOptions(err error, options TraceOptions) error {
	return traceWithOptions(err, options)
}

// traceWithOptions is necessary to match the number of stack frames for
// TraceWithOptions to the other public function such that TraceSkipInternal
// works for all functions.
func traceWithOptions(err error, options TraceOptions) error {
	if err == nil {
		return nil
	}

	if options.Capturer == nil {
		options.Capturer = defaultCapturer
	}

	return &TraceError{
		Data: options.Capturer.Capture(err, options.Skip),
		Err:  err,
	}
}
