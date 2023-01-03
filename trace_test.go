package oops_test

import (
	"errors"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

func TestTrace(t *testing.T) {
	t.Run("folding", func(t *testing.T) {
		require.Nil(t, oops.Trace(nil))

		// NOTE: This isn't the same as the check above (which checks
		// both object == nil and the reflection value:
		// https://github.com/stretchr/testify/blob/master/assert/assertions.go#L519-L539).
		// require.NoError only performs the basic object == nil kind
		// of check (as a typical Go programmer would do).
		//
		// An interface is composed of both the type and value. When
		// oops.Trace returned the pointer *TraceError instead of the
		// interface error the type information was passed along.
		// Unfortunately that meant when passing *TraceError to a
		// function that takes error values (e.g. Chain) the error
		// interface had non-empty type information. This meant that
		// err == nil would return false in such functions because
		// interface{*oops.Trace, nil} is not the same as
		// interface{nil, nil}. This is unfortunate since it means we
		// can't follow the usual paradigm of returning concrete types
		// and accepting interfaces... Or at least we can't and still
		// rely on the widely accepted `err == nil` practice. The
		// options for not checking nil in this way are all deficient
		// in some way (performance, convention).
		//
		// Further reading for the curious:
		//
		// https://dave.cheney.net/2017/08/09/typed-nils-in-go-2
		require.NoError(t, oops.Trace(nil))
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
		terr := func() error {
			return oops.Trace(errors.New("bad stuff"))
		}()

		t.Logf("%+v", terr)

		tnerr := func() error {
			return oops.TraceN(errors.New("bad stuff"), 6)
		}()

		t.Logf("%+v", tnerr)

		require.Equal(t, len(terr.(*oops.TraceError).Data.(oops.Frames)), len(tnerr.(*oops.TraceError).Data.(oops.Frames)))

		t.Run("Trace", func(t *testing.T) {
			require.Error(t, terr)

			te := &oops.TraceError{}
			require.ErrorAs(t, terr, &te)

			t.Logf("%+v", te)

			// NOTE: 3 frames are expected: one for where trace is called
			// in first, one for where first itself is called, and finally
			// one for the test runner.
			require.Len(t, te.Data.(oops.Frames), 3)
		})

		t.Run("TraceN", func(t *testing.T) {
			require.Error(t, tnerr)

			te := &oops.TraceError{}
			require.ErrorAs(t, tnerr, &te)

			t.Logf("%+v", te)

			// NOTE: 3 frames are expected: one for where trace is called
			// in first, one for where first itself is called, and finally
			// one for the test runner.
			require.Len(t, te.Data.(oops.Frames), 3)
		})
	})
}
