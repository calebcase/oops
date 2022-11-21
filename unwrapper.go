package oops

// unwrapper implements the errors' package implied Unwrap interface.
type unwrapper interface {
	Unwrap() error
}
