package oops

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/calebcase/oops/lines"
)

// ShadowError is an error with a hidden error inside.
type ShadowError struct {
	Hidden error `json:"hidden"`
	Err    error `json:"err"`
}

var _ error = &ShadowError{}
var _ unwrapper = &ShadowError{}

func (se *ShadowError) Error() string {
	return fmt.Sprintf("%v", se)
}

func (se *ShadowError) Unwrap() error {
	if se == nil || se.Err == nil || se.Hidden == nil {
		return nil
	}

	if un, ok := se.Err.(interface{ Unwrap() error }); ok {
		return un.Unwrap()
	}

	return nil
}

// Format implements fmt.Format.
func (se *ShadowError) Format(f fmt.State, verb rune) {
	if se == nil || se.Err == nil || se.Hidden == nil {
		return
	}

	flag := ""
	if f.Flag(int('+')) {
		flag = "+"
	}

	if flag == "" {
		fmt.Fprintf(f, "%"+string(verb), se.Err)

		return
	}

	output := lines.Indent(lines.Sprintf("%"+flag+string(verb), se.Err), "··", 1)
	hidden := lines.Indent(lines.Sprintf("%"+flag+string(verb), se.Hidden), "··", 1)

	output = append(output, "··hidden: "+hidden[0])

	if len(hidden) > 1 {
		output = append(output, hidden[1:]...)
	}

	fmt.Fprintf(f, strings.Join(output, "\n"))
}

// MarshalJSON implements json.Marshaler.
func (se *ShadowError) MarshalJSON() (bs []byte, err error) {
	if se == nil || se.Err == nil {
		return []byte("null"), nil
	}

	ebs, err := ErrorMarshalJSON(se.Err)
	if err != nil {
		return nil, err
	}

	output := struct {
		Type string          `json:"type"`
		Err  json.RawMessage `json:"err"`
	}{
		Type: fmt.Sprintf("%T", se.Err),
		Err:  json.RawMessage(ebs),
	}

	return json.Marshal(output)
}

// Shadow hides internal errors with another error.
func Shadow(hidden, err error) error {
	if hidden == nil || err == nil {
		return nil
	}

	return &ShadowError{
		Hidden: hidden,
		Err:    err,
	}
}

// ShadowP replaces hidden with error.
func ShadowP(hidden *error, err error) {
	if hidden == nil {
		return
	}

	*hidden = Shadow(*hidden, err)
}

// ShadowF returns a function that shadows hidden in place.
func ShadowF(hidden *error, err error) func() {
	return func() {
		ShadowP(hidden, err)
	}
}
