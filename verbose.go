package oops

import (
	"fmt"
)

// VerboseError is an error that uses %+v for detailed formatting.
type VerboseError struct {
	Err error
}

// Error implements error.
func (ve *VerboseError) Error() string {
	return fmt.Sprintf("%+v", ve.Err)
}

// Unwrap implements the implied interface for errors.Unwrap.
func (ve *VerboseError) Unwrap() error {
	if ve == nil || ve.Err == nil {
		return nil
	}

	return ve.Err
}

// Verbose wraps an error making its string output verbose ("%+v").
func Verbose(err error) error {
	return &VerboseError{
		Err: err,
	}
}
