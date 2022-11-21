package oops

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Capturer provides a method for capturing trace data.
type Capturer interface {
	Capture(err error, skip int) (data any)
}

// Capture is used to capture trace data. By default it is configured to
// capture stack frames.
var Capture = func(_ error, skip int) (data any) {
	callers := make([]uintptr, 10)
	n := runtime.Callers(skip, callers)
	callers = callers[:n]

	cfs := runtime.CallersFrames(callers)

	frames := make([]runtime.Frame, 0, n)
	for {
		frame, more := cfs.Next()
		if !more {
			break
		}

		frames = append(frames, frame)
	}

	return frames
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

// Error implements errors.Unwraper.
func (te *TraceError) Unwrap() error {
	if te == nil || te.Err == nil {
		return nil
	}

	if uw, ok := te.Err.(unwrapper); ok {
		return uw.Unwrap()
	}

	return nil
}

// Format implements fmt.Format.
func (te *TraceError) Format(f fmt.State, verb rune) {
	if te == nil || te.Err == nil {
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

	output := ErrorIndent(te.Err, "%"+flag+string(verb), "··")

	f.Write([]byte(strings.Join(output, "\n")))

	// TODO: Indent?
	fmt.Fprintf(f, "%"+flag+string(verb), te.Data)
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
// the trace data and 2 levels are skipped.
func Trace(err error) *TraceError {
	return TraceN(err, 2)
}

// TraceN captures a trace and combines it with err. Capturer is use to capture
// the trace data with skipped levels.
func TraceN(err error, skip int) *TraceError {
	return TraceWithOptions(err, TraceOptions{
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
func TraceWithOptions(err error, options TraceOptions) *TraceError {
	if err == nil {
		return nil
	}

	capture := Capture
	if options.Capturer != nil {
		capture = options.Capturer.Capture
	}

	return &TraceError{
		Data: capture(err, options.Skip),
		Err:  err,
	}
}
