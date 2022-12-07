package oops_test

import (
	"errors"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

var ErrShadow = errors.New("shadow")

func shadowFunc() (err error) {
	defer oops.ShadowP(&err, ErrShadow)

	return oops.New("hidden internal error")
}

func TestShadow(t *testing.T) {
	err := shadowFunc()
	require.Error(t, err)

	t.Logf("err: %v\n", err)
	t.Logf("err: %+v\n", err)

	se := err.(*oops.ShadowError)
	require.Equal(t, ErrShadow, se.Err)
}
