package oops_test

import (
	"context"
	"encoding/json"
	"errors"
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
}
