package oops_test

import (
	"errors"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

type CustomError struct {
	Message string
}

func NewCustomError(msg string) *CustomError {
	return &CustomError{
		Message: msg,
	}
}

var _ error = &CustomError{}

func (ce *CustomError) Error() string {
	return ce.Message
}

func TestCustomError(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		cerr := NewCustomError("custom")
		err := oops.New("%w", cerr)

		ce := &CustomError{}
		require.ErrorIs(t, err, cerr)
		require.ErrorAs(t, err, &ce)

		require.Equal(t, cerr, ce)
		require.Equal(t, cerr, errors.Unwrap(errors.Unwrap(err)))
	})

	t.Run("trace", func(t *testing.T) {
		cerr := NewCustomError("custom")
		err := oops.Trace(cerr)

		ce := &CustomError{}
		require.ErrorIs(t, err, cerr)
		require.ErrorAs(t, err, &ce)

		require.Equal(t, cerr, ce)
		require.Equal(t, cerr, errors.Unwrap(err))
	})

	t.Run("namespace", func(t *testing.T) {
		ns := oops.Namespace("test")
		cerr := NewCustomError("custom")
		err := ns.New(cerr)

		ce := &CustomError{}
		require.ErrorIs(t, err, cerr)
		require.ErrorAs(t, err, &ce)

		require.Equal(t, ce, cerr)
		require.Equal(t, ce, errors.Unwrap(err))
	})

	t.Run("shadow", func(t *testing.T) {
		herr := errors.New("hidden")
		cerr := NewCustomError("custom")
		err := oops.Shadow(herr, cerr)

		ce := &CustomError{}
		require.ErrorIs(t, err, cerr)
		require.ErrorAs(t, err, &ce)

		require.Equal(t, ce, cerr)
		require.Equal(t, ce, errors.Unwrap(err))
	})

	t.Run("chain", func(t *testing.T) {
		first := errors.New("first")
		cerr := NewCustomError("custom")
		err := oops.Chain(first, cerr)

		ce := &CustomError{}
		require.ErrorIs(t, err, cerr)
		require.ErrorAs(t, err, &ce)

		require.Equal(t, cerr, ce)
		require.Equal(t, oops.ChainError{cerr}, errors.Unwrap(err))
		require.Equal(t, cerr, errors.Unwrap(errors.Unwrap(err)))
	})
}
