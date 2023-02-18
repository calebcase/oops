package oops

import "fmt"

// Namespace provides a name prefix for new errors.
type Namespace string

var _ Namespacer = Namespace("")

// Wrap returns the error wrapped in the namespace.
func (n Namespace) Wrap(err error) error {
	if err == nil {
		return nil
	}

	return &NamespaceError{
		Name: string(n),
		Err:  err,
	}
}

// WrapP replaces the error with one wrapped in the namespace.
func (n Namespace) WrapP(err *error) {
	if err == nil || *err == nil {
		return
	}

	*err = n.Wrap(*err)
}

func (n Namespace) New(format string, a ...any) error {
	return n.Wrap(TraceN(fmt.Errorf(format, a...), TraceSkipInternal))
}

func (n Namespace) Trace(err error) error {
	return n.Wrap(TraceN(err, TraceSkipInternal))
}

func (n Namespace) TraceN(err error, skip int) error {
	return n.Wrap(TraceN(err, skip))
}

func (n Namespace) TraceWithOptions(err error, options TraceOptions) error {
	return n.Wrap(TraceWithOptions(err, options))
}

func (n Namespace) Chain(errs ...error) error {
	return n.Wrap(Chain(errs...))
}

func (n Namespace) ChainP(err *error, errs ...error) {
	ChainP(err, errs...)
	n.WrapP(err)
}

func (n Namespace) Shadow(hidden, err error) error {
	return n.Wrap(Shadow(hidden, err))
}

func (n Namespace) ShadowP(hidden *error, err error) {
	ShadowP(hidden, err)
	n.WrapP(hidden)
}

func (n Namespace) ShadowF(hidden *error, err error) func() {
	return func() {
		n.ShadowP(hidden, err)
	}
}

type Namespacer interface {
	Wrap(err error) error
	WrapP(err *error)

	New(format string, a ...any) error

	Trace(err error) error
	TraceN(err error, skip int) error
	TraceWithOptions(err error, options TraceOptions) error

	Chain(errs ...error) error
	ChainP(err *error, errs ...error)

	Shadow(hidden, err error) error
	ShadowP(hidden *error, err error)
	ShadowF(hidden *error, err error) func()
}
