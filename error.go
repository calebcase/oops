package oops

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ErrorMarshalJSON uses e's json.Marshaler if it implements one otherwise it
// uses the output from e.Error() (marshalled into a JSON string).
func ErrorMarshalJSON(e error) (bs []byte, err error) {
	if jm, ok := e.(json.Marshaler); ok {
		bs, err = jm.MarshalJSON()
	} else {
		bs, err = json.Marshal(e.Error())
	}

	if err != nil {
		return nil, err
	}

	return bs, nil
}

// ErrorLines returns the error's string in the given format as an array of
// lines.
func ErrorLines(err error, format string) []string {
	return strings.Split(fmt.Sprintf(format, err), "\n")
}

// ErrorIndent returns the error's string prefixed by indent. The first line is
// not indented.
func ErrorIndent(err error, format, indent string) []string {
	lines := ErrorLines(err, format)
	for i, line := range lines {
		if i == 0 {
			continue
		}

		lines[i] = indent + line
	}

	return lines
}
