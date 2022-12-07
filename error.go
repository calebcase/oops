package oops

import (
	"encoding/json"
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
