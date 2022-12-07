package oops_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/calebcase/oops"
	"github.com/stretchr/testify/require"
)

type closer struct {
	open bool
}

func open() (*closer, error) {
	return &closer{true}, nil
}

func (c *closer) Close(ctx context.Context) error {
	c.open = false

	return oops.New("close error")
}

func openclose() (err error) {
	ctx := context.Background()

	f, err := open()
	if err != nil {
		return err
	}
	defer func() {
		err = oops.Chain(err, f.Close(ctx))
	}()

	return oops.New("openclose error")
}

func opencloseP() (err error) {
	ctx := context.Background()

	f, err := open()
	if err != nil {
		return err
	}
	defer func() {
		oops.ChainP(&err, f.Close(ctx))
	}()

	return oops.New("openclose error")
}

func TestChain(t *testing.T) {
	t.Run("Chain", func(t *testing.T) {
		err := openclose()
		werr := errors.Unwrap(err)

		t.Logf("error: %T\n%+v\n", err, err)
		require.Error(t, err)

		t.Logf("wrapped error: %T\n%+v\n", werr, werr)
		require.Error(t, werr)

		bs, jerr := json.MarshalIndent(err, "", "  ")
		require.NoError(t, jerr)
		t.Log(string(bs))
	})

	t.Run("ChainP", func(t *testing.T) {
		err := opencloseP()
		werr := errors.Unwrap(err)

		t.Logf("error: %T\n%+v\n", err, err)
		require.Error(t, err)

		t.Logf("wrapped error: %T\n%+v\n", werr, werr)
		require.Error(t, werr)
	})

	// Chain should "fold" away nils and remove itself if the chain only
	// contains one element.
	t.Run("folding", func(t *testing.T) {
		type TC struct {
			Errs   []error
			Expect error
		}

		e0 := errors.New("0")
		e1 := errors.New("1")
		e2 := errors.New("2")

		tcs := []TC{
			{
				Errs:   nil,
				Expect: nil,
			},
			{
				Errs:   []error{},
				Expect: nil,
			},
			{
				Errs:   []error{nil},
				Expect: nil,
			},
			{
				Errs:   []error{e0},
				Expect: e0,
			},
			{
				Errs:   []error{nil, nil},
				Expect: nil,
			},
			{
				Errs:   []error{nil, e1},
				Expect: e1,
			},
			{
				Errs:   []error{e0, nil},
				Expect: e0,
			},
			{
				Errs:   []error{e0, e1},
				Expect: oops.ChainError{e0, e1},
			},
			{
				Errs:   []error{nil, nil, nil},
				Expect: nil,
			},
			{
				Errs:   []error{nil, nil, e2},
				Expect: e2,
			},
			{
				Errs:   []error{nil, e1, nil},
				Expect: e1,
			},
			{
				Errs:   []error{nil, e1, e2},
				Expect: oops.ChainError{e1, e2},
			},
			{
				Errs:   []error{e0, nil, nil},
				Expect: e0,
			},
			{
				Errs:   []error{e0, nil, e2},
				Expect: oops.ChainError{e0, e2},
			},
			{
				Errs:   []error{e0, e1, nil},
				Expect: oops.ChainError{e0, e1},
			},
			{
				Errs:   []error{e0, e1, e2},
				Expect: oops.ChainError{e0, e1, e2},
			},
		}

		for _, tc := range tcs {
			t.Run(fmt.Sprintf("%s", tc.Errs), func(t *testing.T) {
				require.Equal(t, tc.Expect, oops.Chain(tc.Errs...))
			})
		}
	})
}
