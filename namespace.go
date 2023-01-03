package oops

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/calebcase/oops/lines"
)

// Namespace provides a name prefix for new errors.
type Namespace string

// New returns the error wrapped in the namespace.
func (n Namespace) New(err error) error {
	if err == nil {
		return nil
	}

	return &NamespaceError{
		Name: string(n),
		Err:  err,
	}
}

// NewP replaces err with a namespaced error.
func (n Namespace) NewP(err *error) {
	if err == nil || *err == nil {
		return
	}

	*err = n.New(*err)
}

// NamespaceError prefixes errors with the given namespace.
type NamespaceError struct {
	Name string
	Err  error
}

var _ error = &NamespaceError{}
var _ unwrapper = &NamespaceError{}

func (ne *NamespaceError) Error() string {
	if ne == nil || ne.Err == nil {
		return ""
	}

	return ne.Name + ": " + ne.Err.Error()
}

func (ne *NamespaceError) Unwrap() error {
	if ne == nil {
		return nil
	}

	return ne.Err
}

// Format implements fmt.Format.
func (ne *NamespaceError) Format(f fmt.State, verb rune) {
	if ne == nil || ne.Err == nil {
		fmt.Fprintf(f, "<nil>")

		return
	}

	flag := ""
	if f.Flag(int('+')) {
		flag = "+"
	}

	if flag == "" {
		fmt.Fprintf(f, "%s: %"+string(verb), ne.Name, ne.Err)

		return
	}

	output := []string{}
	ls := lines.Sprintf("%"+flag+string(verb), ne.Err)

	output = append(output, fmt.Sprintf("%s: %s", ne.Name, ls[0]))

	if len(ls) > 1 {
		output = append(output, ls[1:]...)
	}

	f.Write([]byte(strings.Join(output, "\n")))
}

// MarshalJSON implements json.Marshaler.
func (ne *NamespaceError) MarshalJSON() (bs []byte, err error) {
	if ne == nil || ne.Err == nil {
		return []byte("null"), nil
	}

	ebs, err := ErrorMarshalJSON(ne.Err)
	if err != nil {
		return nil, err
	}

	output := struct {
		Type string          `json:"type"`
		Err  json.RawMessage `json:"err"`
	}{
		Type: fmt.Sprintf("%T", ne.Err),
		Err:  json.RawMessage(ebs),
	}

	return json.Marshal(output)
}
