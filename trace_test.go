package oops_test

import (
	"errors"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

func TestTrace(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		first := func() error {
			return oops.Trace(errors.New("bad stuff"))
		}

		second := func() error {
			return first()
		}

		ferr := first()
		require.Error(t, ferr)

		te := &oops.TraceError{}
		require.ErrorAs(t, ferr, &te)

		t.Logf("%+v", te)
		require.Len(t, te.Data.(oops.Frames), 8)

		serr := second()
		require.Error(t, serr)
		require.ErrorAs(t, serr, &te)

		t.Logf("%+v", te)
		require.Len(t, te.Data.(oops.Frames), 9)
	})

	t.Run("skip_0", func(t *testing.T) {
		first := func() error {
			return oops.TraceN(errors.New("bad stuff"), 0)
		}

		ferr := first()
		require.Error(t, ferr)

		te := &oops.TraceError{}
		require.ErrorAs(t, ferr, &te)

		t.Logf("%+v", te)
		require.Len(t, te.Data.(oops.Frames), 9)
	})

	t.Run("skip_1", func(t *testing.T) {
		first := func() error {
			return oops.TraceN(errors.New("bad stuff"), 1)
		}

		ferr := first()
		require.Error(t, ferr)

		te := &oops.TraceError{}
		require.ErrorAs(t, ferr, &te)

		t.Logf("%+v", te)
		require.Len(t, te.Data.(oops.Frames), 8)
	})

	t.Run("skip_internal", func(t *testing.T) {
		first := func() error {
			return oops.TraceN(errors.New("bad stuff"), 6)
		}

		ferr := first()
		require.Error(t, ferr)

		te := &oops.TraceError{}
		require.ErrorAs(t, ferr, &te)

		t.Logf("%+v", te)

		// NOTE: 3 frames are expected: one for where trace is called
		// in first, one for where first itself is called, and finally
		// one for the test runner.
		require.Len(t, te.Data.(oops.Frames), 3)
	})
}
