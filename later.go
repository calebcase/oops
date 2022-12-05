package oops

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type LaterError struct {
	Errorer interface{}
	Args    []interface{}
	Err     *error
}

var _ error = &LaterError{}
var _ unwrapper = &LaterError{}

func (le *LaterError) Eval() (err error) {
	if le == nil || le.Errorer == nil {
		return nil
	}

	if le.Err != nil {
		return *le.Err
	}

	t := reflect.TypeOf(le.Errorer)

	if t.Kind() != reflect.Func {
		return New("errorer is not a function: %T", le.Errorer)
	}

	if t.NumOut() == 0 {
		return New("errorer must return at least an error: %T", le.Errorer)
	}

	fn := reflect.ValueOf(le.Errorer)

	args := make([]reflect.Value, 0, t.NumIn())
	for _, arg := range le.Args {
		args = append(args, reflect.ValueOf(arg))
	}

	returns := fn.Call(args)
	last := returns[len(returns)-1]

	var ok bool
	if err, ok = last.Interface().(error); !ok {
		return New("errorer's last arg is not an error: %T", le.Errorer)
	}
	le.Err = &err

	return err
}

func (le *LaterError) Error() string {
	if le == nil {
		return ""
	}

	return le.Eval().Error()
}

func (le *LaterError) Unwrap() error {
	if le == nil {
		return nil
	}

	return le.Eval()
}

// Format implements fmt.Format.
func (le *LaterError) Format(f fmt.State, verb rune) {
	le.Eval()

	if le == nil || le.Err == nil || *le.Err == nil {
		return
	}

	flag := ""
	if f.Flag(int('+')) {
		flag = "+"
	}

	if flag == "" {
		fmt.Fprintf(f, "%"+string(verb), le.Err)

		return
	}

	lines := ErrorLines(*le.Err, "%"+flag+string(verb))

	fmt.Fprintf(f, "later: %s\n", lines[0])

	if len(lines) > 1 {
		f.Write([]byte(strings.Join(lines[1:], "\n")))
	}
}

// MarshalJSON implements json.Marshaler.
func (le *LaterError) MarshalJSON() (bs []byte, err error) {
	le.Eval()

	if le == nil || le.Err == nil || *le.Err == nil {
		return []byte("null"), nil
	}

	ebs, err := ErrorMarshalJSON(*le.Err)
	if err != nil {
		return nil, err
	}

	output := struct {
		Type string          `json:"type"`
		Err  json.RawMessage `json:"err"`
	}{
		Type: fmt.Sprintf("%T", le.Err),
		Err:  json.RawMessage(ebs),
	}

	return json.Marshal(output)
}

// Later will call errorer when later is evaluated. Errorer will be called with
// the provided args and the last returned value will be used as the error.
// errorer must return an error as its last return value.
func Later(errorer interface{}, args ...interface{}) error {
	return TraceN(&LaterError{
		Errorer: errorer,
		Args:    args,
	}, 3)
}
