package oops

import (
	"fmt"
)

// New returns a new error from fmt.Errorf with a stack trace.
func New(format string, a ...any) error {
	return TraceN(fmt.Errorf(format, a...), 3)
}
